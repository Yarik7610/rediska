package commands

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) configGet(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "CONFIG GET command must have only 1 arg"}
	}

	arg := args[0]
	value := []string{arg}
	switch arg {
	case "host":
		value = append(value, c.args.Host)
	case "port":
		value = append(value, strconv.Itoa(c.args.Port))
	case "replicaof":
		var replicaOfString string
		if c.args.ReplicaOf != nil {
			replicaOfString = c.args.ReplicaOf.String()
		}
		value = append(value, replicaOfString)
	case "dir":
		value = append(value, c.args.DBDir)
	case "dbfilename":
		value = append(value, c.args.DBFilename)
	default:
		return resp.SimpleError{Value: fmt.Sprintf("CONFIG GET command unknown arg: %s", arg)}
	}

	return resp.CreateBulkStringArray(value...)
}
