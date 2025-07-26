package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) llen(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "LLEN command must have 1 arg"}
	}

	key := args[0]
	if _, ok := c.storage.StringStorage.Get(key); ok {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	len := c.storage.ListStorage.Llen(key)
	return resp.Integer{Value: len}
}
