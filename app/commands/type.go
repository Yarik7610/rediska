package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) valuetype(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "TYPE command must have only 1 arg"}
	}

	key := args[0]

	if c.storage.StringStorage().Has(key) {
		return resp.SimpleString{Value: memory.TYPE_STRING}
	}
	if c.storage.ListStorage().Has(key) {
		return resp.SimpleString{Value: memory.TYPE_LIST}
	}
	if c.storage.StreamStorage().Has(key) {
		return resp.SimpleString{Value: memory.TYPE_STREAM}
	}

	return resp.SimpleString{Value: memory.TYPE_NONE}
}
