package state

import (
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
