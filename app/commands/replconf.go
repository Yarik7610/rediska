package commands

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

func (c *controller) replconf(args []string, conn net.Conn) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "REPLCONF command error: only 2 more arguments supported"}
	}

	addr := utils.GetRemoteAddr(conn)

	secondCommand := args[0]
	arg := args[1]
	switch replicationController := c.replicationController.(type) {
	case replication.MasterController:
		switch strings.ToLower(secondCommand) {
		case "listening-port":
			replicationController.AddReplicaConn(conn)
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
			replicationController.SendAck(addr, ackOffset)
			return nil
		default:
			return resp.SimpleError{Value: fmt.Sprintf("REPLCONF master unsupported second command: %s", secondCommand)}
		}
	case replication.ReplicaController:
		switch strings.ToUpper(secondCommand) {
		case "GETACK":
			if arg != "*" {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF GETACK replica unsupported argument: %s", arg)}
			}
			if replicationController.GetMasterConn() != conn {
				return resp.SimpleError{Value: "REPLCONF GETACK * can be send only by master"}
			}

			// No syncing with propagated write commands from master
			// Ideally, we should wait here, until replica won't have any write commands from master in processing
			// Or we need to compare that Ack.offset >= master.MasterReplOffset and only than push Ack to ackedReplicas (wait.go)
			response := resp.CreateBulkStringArray("REPLCONF", "ACK", strconv.Itoa(replicationController.Info().MasterReplOffset))
			if err := utils.WriteCommand(response, conn); err != nil {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF GETACK * write to master error: %v", err)}
			}
			return nil
		}
		return resp.SimpleError{Value: fmt.Sprintf("REPLCONF isn't supported for replica: %s", secondCommand)}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("REPLCONF command detected unknown type assertion: %T", replicationController)}
	}
}
