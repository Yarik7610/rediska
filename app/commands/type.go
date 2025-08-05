package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) valuetype(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "TYPE command must have only 1 arg"}
	}

	key := args[0]

	return resp.SimpleString{Value: c.storage.Type(key)}
}
