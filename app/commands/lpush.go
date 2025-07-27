package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) lpush(args, commandAndArgs []string) resp.Value {
	if len(args) < 2 {
		return resp.SimpleError{Value: "LPUSH command must have at least 2 args"}
	}

	key := args[0]
	values := args[1:]

	if _, ok := c.storage.StringStorage.Get(key); ok {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	len := c.storage.ListStorage.Lpush(key, values...)
	c.propagateWriteCommand(commandAndArgs)
	return resp.Integer{Value: len}
}
