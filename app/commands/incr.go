package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) incr(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "INCR command must have only 1 arg"}
	}

	key := args[0]
	if c.storage.KeyExistsWithOtherType(key, memory.TYPE_STRING) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	incremented, err := c.storage.StringStorage().Incr(key)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("ERR %s", err)}
	}

	return resp.Integer{Value: incremented}
}
