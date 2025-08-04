package commands

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) rpop(args, commandAndArgs []string) resp.Value {
	if len(args) < 1 {
		return resp.SimpleError{Value: "RPOP command must have at least 1 arg"}
	}

	key := args[0]
	if c.storage.KeyExistsWithOtherType(key, memory.TYPE_LIST) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	count := 1
	if len(args) > 1 {
		var err error
		count, err = strconv.Atoi(args[1])
		if err != nil {
			return resp.SimpleError{Value: fmt.Sprintf("RPOP command count argument atoi error: %v", err)}
		}
	}

	poppedValues := c.storage.ListStorage().Rpop(key, count)
	c.propagateWriteCommand(commandAndArgs)

	switch len(poppedValues) {
	case 0:
		return resp.BulkString{Value: nil}
	case 1:
		return resp.BulkString{Value: &poppedValues[0]}
	default:
		return resp.CreateBulkStringArray(poppedValues...)
	}
}
