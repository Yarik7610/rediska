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

	replID := args[0]
	replOffset := args[1]

	switch r := c.replication.(type) {
	case replication.Master:
		if replID == "?" && replOffset == "-1" {
			// return resp.CreateArray()
		}
		return resp.SimpleError{Value: fmt.Sprintf("psync master unsupported replication id: %s and replication offset: %s", replID, replOffset)}
	case replication.Replica:
		return resp.SimpleError{Value: "psync isn't supported for replica"}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("psync command detected unknown type assertion: %T", r)}
	}
}
