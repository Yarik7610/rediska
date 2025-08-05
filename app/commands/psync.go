package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

func (c *controller) psync(args []string, conn net.Conn) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "PSYNC command error: only 2 argument supported"}
	}

	requestedReplID := args[0]
	requestedReplOffset := args[1]

	switch replicationController := c.replicationController.(type) {
	case replication.MasterController:
		if requestedReplID == "?" && requestedReplOffset == "-1" {
			if !replicationController.IsReplica(conn) {
				return resp.SimpleError{Value: "PSYNC command error: failed to send FULLRESYNC, because no such replica exists"}
			}

			response := "FULLRESYNC" + " " + replicationController.Info().MasterReplID + " " + "0"
			if err := utils.WriteCommand(resp.SimpleString{Value: response}, conn); err != nil {
				return resp.SimpleError{Value: fmt.Sprintf("PSYNC command error: failed to send FULLRESYNC: %v", err)}
			}
			replicationController.SendRDBFile(conn)
			return nil
		}
		return resp.SimpleError{Value: fmt.Sprintf("PSYNC master unsupported replication id: %s and replication offset: %s", requestedReplID, requestedReplOffset)}
	case replication.ReplicaController:
		return resp.SimpleError{Value: "PSYNC isn't supported for replica"}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("PSYNC command detected unknown type assertion: %T", replicationController)}
	}
}
