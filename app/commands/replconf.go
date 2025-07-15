package commands

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) replconf(args []string) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "REPLCONF command error: only 2 argument supported"}
	}

	section := args[0]
	arg := args[1]
	switch section {
	case "listening-port":
		atoiPort, err := strconv.Atoi(arg)
		if err != nil {
			return resp.SimpleError{Value: fmt.Sprintf("REPLCONF wrong 'listening-port' argument: %s", arg)}
		}
		if atoiPort != c.args.ReplicaOf.Port {
			return resp.SimpleError{Value: fmt.Sprintf("REPLCONF 'listening-port' argument doesn't match replica of port: %d", c.args.ReplicaOf.Port)}
		}
		//TODO
		return resp.SimpleString{Value: "OK"}
	case "capa":
		if arg != "psync2" {
			return resp.SimpleError{Value: fmt.Sprintf("REPLCONF unsupported argument for 'capa': %s", arg)}
		}
		return resp.SimpleString{Value: "OK"}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("REPLCONF unsupported section: %s", section)}
	}
}
