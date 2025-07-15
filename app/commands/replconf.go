package commands

import (
	"fmt"
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) replconf(args []string, conn net.Conn) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "REPLCONF command error: only 2 argument supported"}
	}

	section := args[0]
	arg := args[1]

	switch r := c.replication.(type) {
	case replication.Master:
		switch section {
		case "listening-port":
			addr := conn.RemoteAddr().String()
			r.AddReplicaConn(addr, conn)
			return resp.SimpleString{Value: "OK"}
		default:
			return resp.SimpleError{Value: fmt.Sprintf("REPLCONF unsupported section: %s", section)}
		}
	case replication.Replica:
		switch section {
		case "listening-port":
			atoiPort, err := strconv.Atoi(arg)
			if err != nil {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF wrong 'listening-port' argument: %s", arg)}
			}
			if atoiPort != c.args.Port {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF 'listening-port' argument doesn't match replica port: %d", c.args.Port)}
			}
			return resp.SimpleString{Value: "OK"}
		case "capa":
			if arg != "psync2" {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF unsupported argument for 'capa': %s", arg)}
			}
			return resp.SimpleString{Value: "OK"}
		default:
			return resp.SimpleError{Value: fmt.Sprintf("REPLCONF unsupported section: %s", section)}
		}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("REPLCONF command detected unknown type assertion: %T", r)}
	}
}
