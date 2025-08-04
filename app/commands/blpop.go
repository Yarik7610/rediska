package commands

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) blpop(args, commandAndArgs []string) resp.Value {
	if len(args) < 2 {
		return resp.SimpleError{Value: "BLPOP command must have at least 2 args"}
	}

	key := args[0]
	if c.storage.KeyExistsWithOtherType(key, memory.TYPE_LIST) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	timeoutS, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("BLPOP command timeout (S) argument parseFloat error: %v", err)}
	}
	poppedValue := c.storage.ListStorage().Blpop(key, timeoutS)

	c.propagateWriteCommand(commandAndArgs)

	if poppedValue == nil {
		return resp.BulkString{Value: nil}
	}
	return resp.CreateBulkStringArray(key, *poppedValue)
}
