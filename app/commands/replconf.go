package commands

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) replconf(args []string, conn net.Conn) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "REPLCONF command error: only 2 more arguments supported"}
	}

	secondCommand := args[0]
	arg := args[1]

	switch r := c.replication.(type) {
	case replication.Master:
		switch secondCommand {
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
			return resp.SimpleError{Value: fmt.Sprintf("REPLCONF master unsupported section: %s", secondCommand)}
		}
	case replication.Replica:
		switch strings.ToUpper(secondCommand) {
		case "GETACK":
			if arg != "*" {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF GETACK unsupported argument: %s", arg)}
			}
			if r.GetMasterConn() != conn {
				return resp.SimpleError{Value: "REPLCONF GETACK * can be send only by master"}
			}

			response := resp.CreateBulkStringArray("REPLCONF", "ACK", strconv.Itoa(r.Info().MasterReplOffset))
			if err := c.Write(response, conn); err != nil {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF GETACK response error: %v", err)}
			}
		}
		return resp.SimpleError{Value: fmt.Sprintf("REPLCONF isn't supported for replica: %s", secondCommand)}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("REPLCONF command detected unknown type assertion: %T", r)}
	}
}
