package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) get(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "GET command must have only 1 arg"}
	}

	key := args[0]

	_, ok := c.storage.ListStorage.Get(key)
	if ok {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	got, ok := c.storage.StringStorage.Get(key)
	if !ok {
		return resp.BulkString{Value: nil}
	}

	return resp.BulkString{Value: &got.Value}
}
