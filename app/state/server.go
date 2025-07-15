package state

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/persistence/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func SpawnServer(args *config.Args) Server {
	if args.ReplicaOf == nil {
		return newMasterServer(args)
	} else {
		return newReplicaServer(args)
	}
}

type Server interface {
	Start()
}

type baseServer struct {
	Storage           *memory.Storage
	RESPController    *resp.Controller
	Args              *config.Args
	CommandController *commands.Controller
	ReplicationInfo   *replication.Info
}

func newBaseServer(args *config.Args) *baseServer {
	return &baseServer{
		Storage:        memory.NewStorage(),
		RESPController: resp.NewController(),
		Args:           args,
	}
}

func (s *baseServer) initStorage() {
	if s.Args.DBDir == "" || s.Args.DBFilename == "" {
		return
	}

	if !rdb.IsFileExists(s.Args.DBDir, s.Args.DBFilename) {
		return
	}

	b, err := rdb.ReadRDB(s.Args.DBDir, s.Args.DBFilename)
	if err != nil {
		log.Printf("Skip RDB storage seed, RDB file read error: %v\n", err)
		return
	}

	items, err := rdb.Decode(b)
	if err != nil {
		log.Printf("Skip RDB storage seed, RDB decode error: %v\n", err)
		return
	}
	s.putRDBItemsIntoStorage(items)
}

func (s *baseServer) putRDBItemsIntoStorage(items map[string]memory.Item) {
	for key, item := range items {
		if memory.ItemHasExpiration(&item) {
			if !memory.ItemExpired(&item) {
				s.Storage.SetWithExpiry(key, item.Value, item.Expires)
			}
		} else {
			s.Storage.Set(key, item.Value)
		}
	}
}

func (s *baseServer) acceptClientConnections() {
	address := fmt.Sprintf("%s:%d", s.Args.Host, s.Args.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to bind to address: %s\n", address)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go s.handleClient(conn)
	}
}

func (s *baseServer) handleClient(conn net.Conn) {
	defer conn.Close()

	for {
		b := make([]byte, 1024)

		n, err := conn.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Printf("Connnection read error: %v", err)
			return
		}

		value, err := s.RESPController.Decode(b[:n])
		if err != nil {
			fmt.Fprintf(conn, "RESP controller decode error: %v\n", err)
			continue
		}

		s.CommandController.HandleCommand(value, conn)
	}
}

func (s *baseServer) startExpiredKeysCleanup() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			s.Storage.CleanExpiredKeys()
		}
	}()
}
