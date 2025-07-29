package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) rpush(args, commandAndArgs []string) resp.Value {
	if len(args) < 2 {
		return resp.SimpleError{Value: "RPUSH command must have at least 2 args"}
	}

	key := args[0]
	if c.storage.KeyExistsWithOtherType(key, memory.TYPE_LIST) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	values := args[1:]
	len := c.storage.ListStorage().Rpush(key, values...)
	go c.propagateWriteCommand(commandAndArgs)
	return resp.Integer{Value: len}
}
