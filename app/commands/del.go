package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) del(args, commandAndArgs []string) resp.Value {
	for _, key := range args {
		c.storage.Del(key)
	}

	c.propagateWriteCommand(commandAndArgs)
	return resp.SimpleString{Value: "OK"}
}
