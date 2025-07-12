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

type Server interface {
	Start()
	IsMaster() bool

	DecodeRESP(data []byte) (resp.Value, error)
	HandleCommand(value resp.Value, conn net.Conn)
}

type BaseServer struct {
	Listener          net.Listener
	Storage           *memory.Storage
	RESPController    *resp.Controller
	Args              *config.Args
	CommandController *commands.Controller
	ReplicationInfo   *replication.Info
}

func NewBaseServer(args *config.Args, listener net.Listener) *BaseServer {
	return &BaseServer{
		Listener:       listener,
		Storage:        memory.NewStorage(),
		RESPController: resp.NewController(),
		Args:           args,
	}
}

func SpawnServer(args *config.Args, listener net.Listener) Server {
	if args.ReplicaOf == nil {
		return NewMasterServer(args, listener)
	} else {
		return NewReplicaServer(args, listener)
	}
}

func (s *BaseServer) DecodeRESP(b []byte) (resp.Value, error) {
	value, err := s.RESPController.Decode(b)
	return value, err
}

func (s *BaseServer) HandleCommand(value resp.Value, conn net.Conn) {
	s.CommandController.HandleCommand(value, conn)
}

func (s *BaseServer) initStorage() {
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

func (s *BaseServer) putRDBItemsIntoStorage(items map[string]memory.Item) {
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

func (s *BaseServer) acceptConnections() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go s.handleClient(conn)
	}
}

func (s *BaseServer) handleClient(conn net.Conn) {
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

		value, err := s.DecodeRESP(b[:n])
		if err != nil {
			fmt.Fprintf(conn, "RESP controller decode error: %v\n", err)
			continue
		}

		s.HandleCommand(value, conn)
	}
}

func (s *BaseServer) startExpiredKeysCleanup() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			s.Storage.CleanExpiredKeys()
		}
	}()
}
