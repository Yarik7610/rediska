package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/state"
)

func Keys(args []string, server *state.Server) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "KEYS command error: only 1 argument supported"}
	}

	pattern := args[0]
	if pattern != "*" {
		return resp.SimpleError{Value: "KEYS command error: only '*' pattern supported"}
	}

	keys := server.Storage.GetKeys()
	var value []resp.Value
	for _, key := range keys {
		value = append(value, resp.BulkString{Value: &key})
	}
	return resp.Array{Value: value}
}
