package state

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/db/memory"
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

func (s *Server) StartExpiredKeysCleanup() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			s.Storage.CleanExpiredKeys()
		}
	}()
}
