package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/state"
)

func ConfigGet(args []string, server *state.Server) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "CONFIG GET command must have only 1 arg"}
	}

	arg := args[0]
	switch arg {
	case "dir":
		return resp.Array{Value: []resp.Value{
			resp.BulkString{Value: &arg},
			resp.BulkString{Value: &server.Args.DBDir},
		}}
	case "dbfilename":
		return resp.Array{Value: []resp.Value{
			resp.BulkString{Value: &arg},
			resp.BulkString{Value: &server.Args.DBFilename},
		}}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("CONFIG GET command unknown arg: %s", arg)}
	}
}
