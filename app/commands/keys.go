package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) keys(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "KEYS command error: only 1 argument supported"}
	}

	pattern := args[0]
	if pattern != "*" {
		return resp.SimpleError{Value: "KEYS command error: only '*' pattern supported"}
	}

	keys := c.storage.GetKeys()
	return resp.CreateArray(keys...)
}
