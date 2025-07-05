package state

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Server struct {
	Storage        *db.Storage
	RESPController *resp.Controller
}

func NewServer() *Server {
	return &Server{
		Storage:        db.NewStorage(),
		RESPController: resp.NewController(),
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
