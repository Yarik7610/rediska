package commands

import (
	"fmt"
	"net"

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
		case "capa":
			if arg != "psync2" {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF unsupported argument for 'capa': %s", arg)}
			}
			return resp.SimpleString{Value: "OK"}
		default:
			return resp.SimpleError{Value: fmt.Sprintf("REPLCONF master unsupported section: %s", section)}
		}
	case replication.Replica:
		return resp.SimpleError{Value: fmt.Sprintf("REPLCONF replica unsupported section: %s", section)}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("REPLCONF command detected unknown type assertion: %T", r)}
	}
}
