package state

import (
	"log"
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/persistence/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Server interface {
	Start()
	IsMaster() bool

	DecodeRESP(data []byte) (resp.Value, error)
	HandleCommand(value resp.Value, conn net.Conn)
}

type BaseServer struct {
	Storage           *memory.Storage
	RESPController    *resp.Controller
	CommandController *commands.Controller
	Args              *config.Args
}

func NewBaseServer(args *config.Args) *BaseServer {
	return &BaseServer{
		Storage:        memory.NewStorage(),
		RESPController: resp.NewController(),
		Args:           args,
	}
}

func SpawnServer(args *config.Args) Server {
	if args.ReplicaOf == nil {
		return NewMasterServer(args)
	} else {
		return NewReplicaServer(args)
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

func (s *BaseServer) startExpiredKeysCleanup() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			s.Storage.CleanExpiredKeys()
		}
	}()
}
