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
	value := []resp.Value{
		resp.BulkString{Value: &arg},
	}
	switch arg {
	case "host":
		value = append(value, resp.BulkString{Value: &c.args.Host})
	case "port":
		itoaPort := strconv.Itoa(c.args.Port)
		value = append(value, resp.BulkString{Value: &itoaPort})
	case "replicaof":
		var replicaOfString *string
		if c.args.ReplicaOf != nil {
			str := c.args.ReplicaOf.String()
			replicaOfString = &str
		}
		value = append(value, resp.BulkString{Value: replicaOfString})
	case "dir":
		value = append(value, resp.BulkString{Value: &c.args.DBDir})
	case "dbfilename":
		value = append(value, resp.BulkString{Value: &c.args.DBFilename})
	default:
		return resp.SimpleError{Value: fmt.Sprintf("CONFIG GET command unknown arg: %s", arg)}
	}

	return resp.Array{Value: value}
}
