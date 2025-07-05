package state

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Server struct {
	Storage        *db.Storage
	RESPController *resp.Controller
	Args           *config.Args
}

func NewServer(args *config.Args) *Server {
	return &Server{
		Storage:        db.NewStorage(),
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
