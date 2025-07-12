package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) get(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "GET command must have only 1 arg"}
	}

	got, ok := c.storage.Get(args[0])
	if !ok {
		return resp.BulkString{Value: nil}
	}

	return resp.BulkString{Value: &got.Value}
}
