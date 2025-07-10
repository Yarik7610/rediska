package state

import (
	"log"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/persistence/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Server struct {
	Storage        *memory.Storage
	RESPController *resp.Controller
	Args           *config.Args
}

func NewServer(args *config.Args) *Server {
	return &Server{
		Storage:        memory.NewStorage(),
		RESPController: resp.NewController(),
		Args:           args,
	}
}

func (s *Server) InitStorage() {
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

	err = rdb.Decode(b)
	if err != nil {
		log.Printf("Skip RDB storage seed, RDB decode error: %v\n", err)
		return
	}
}

func (s *Server) StartExpiredKeysCleanup() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			s.Storage.CleanExpiredKeys()
		}
	}()
}
