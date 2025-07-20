package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) psync(args []string, conn net.Conn) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "psync command error: only 2 argument supported"}
	}

	requestedReplID := args[0]
	requestedReplOffset := args[1]

	switch r := c.replication.(type) {
	case replication.Master:
		if requestedReplID == "?" && requestedReplOffset == "-1" {
			response := "FULLRESYNC" + " " + r.Info().MasterReplID + " " + "0"
			if err := c.Write(resp.SimpleString{Value: response}, conn); err != nil {
				return resp.SimpleError{Value: fmt.Sprintf("failed to send FULLRESYNC: %v", err)}
			}
			r.SendRDBFile(conn)
			return nil
		}
		return resp.SimpleError{Value: fmt.Sprintf("psync master unsupported replication id: %s and replication offset: %s", requestedReplID, requestedReplOffset)}
	case replication.Replica:
		return resp.SimpleError{Value: "psync isn't supported for replica"}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("psync command detected unknown type assertion: %T", r)}
	}
}
