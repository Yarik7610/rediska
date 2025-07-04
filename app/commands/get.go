package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/state"
)

func Get(args []string, server *state.Server) resp.Value {
	if len(args) > 1 {
		return resp.SimpleError{Value: "GET command must have only 1 arg"}
	}

	got, ok := server.Storage.Get(args[0])
	if !ok {
		return resp.BulkString{Value: nil}
	}

	return resp.BulkString{Value: &got.Value}
}
