package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) info(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "INFO command error: only 1 argument supported"}
	}

	section := args[0]
	switch section {
	case "replication":
		info, err := replicationInfo(c)
		if err != nil {
			return resp.SimpleError{Value: fmt.Sprintf("INFO replication error: %v", err)}
		}
		return resp.BulkString{Value: &info}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("INFO unsupported section: %s", section)}
	}
}

func replicationInfo(c *Controller) (string, error) {
	if c.serverIsMaster {
		return "role:master", nil
	} else {
		return "", fmt.Errorf("replica replication info isn't supported")
	}
}
