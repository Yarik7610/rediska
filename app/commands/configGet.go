package commands

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) configGet(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "CONFIG GET command must have only 1 arg"}
	}

	arg := args[0]
	switch arg {
	case "host":
		return resp.Array{Value: []resp.Value{
			resp.BulkString{Value: &arg},
			resp.BulkString{Value: &c.args.Host},
		}}
	case "port":
		itoaPort := strconv.Itoa(c.args.Port)
		return resp.Array{Value: []resp.Value{
			resp.BulkString{Value: &arg},
			resp.BulkString{Value: &itoaPort},
		}}
	case "replicaof":
		replicaOfString := c.args.ReplicaOf.String()
		return resp.Array{Value: []resp.Value{
			resp.BulkString{Value: &arg},
			resp.BulkString{Value: &replicaOfString},
		}}
	case "dir":
		return resp.Array{Value: []resp.Value{
			resp.BulkString{Value: &arg},
			resp.BulkString{Value: &c.args.DBDir},
		}}
	case "dbfilename":
		return resp.Array{Value: []resp.Value{
			resp.BulkString{Value: &arg},
			resp.BulkString{Value: &c.args.DBFilename},
		}}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("CONFIG GET command unknown arg: %s", arg)}
	}
}
