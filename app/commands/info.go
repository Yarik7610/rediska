package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) info(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "INFO command error: only 1 argument supported"}
	}

	section := args[0]
	switch section {
	case "replication":
		replicationInfo := c.replicationController.Info().String()
		return resp.BulkString{Value: &replicationInfo}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("INFO unsupported section: %s", section)}
	}
}
