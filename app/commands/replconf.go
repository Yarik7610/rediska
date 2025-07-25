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

	addr := conn.RemoteAddr().String()

	secondCommand := args[0]
	arg := args[1]
	switch rt := c.replication.(type) {
	case replication.Master:
		switch strings.ToLower(secondCommand) {
		case "listening-port":
			rt.AddReplicaConn(addr, conn)
			return resp.SimpleString{Value: "OK"}
		case "capa":
			if arg != "psync2" {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF capa unsupported argument: %s", arg)}
			}
			return resp.SimpleString{Value: "OK"}
		case "ack":
			ackOffset, err := strconv.Atoi(arg)
			if err != nil {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF ACK master offset atoi error: %s", secondCommand)}
			}
			rt.SendAck(addr, ackOffset)
			return nil
		default:
			return resp.SimpleError{Value: fmt.Sprintf("REPLCONF master unsupported second command: %s", secondCommand)}
		}
	case replication.Replica:
		switch strings.ToUpper(secondCommand) {
		case "GETACK":
			if arg != "*" {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF GETACK replica unsupported argument: %s", arg)}
			}
			if rt.GetMasterConn() != conn {
				return resp.SimpleError{Value: "REPLCONF GETACK * can be send only by master"}
			}

			response := resp.CreateBulkStringArray("REPLCONF", "ACK", strconv.Itoa(rt.Info().MasterReplOffset))
			if err := c.Write(response, conn); err != nil {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF GETACK * write to master error: %v", err)}
			}
			return nil
		}
		return resp.SimpleError{Value: fmt.Sprintf("REPLCONF isn't supported for replica: %s", secondCommand)}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("REPLCONF command detected unknown type assertion: %T", rt)}
	}
}
