package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) valuetype(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "TYPE command must have only 1 arg"}
	}

	key := args[0]

	if _, ok := c.storage.StringStorage.Get(key); ok {
		return resp.SimpleString{Value: "string"}
	}
	if _, ok := c.storage.ListStorage.Get(key); ok {
		return resp.SimpleString{Value: "list"}
	}

	return resp.SimpleString{Value: "none"}
}
