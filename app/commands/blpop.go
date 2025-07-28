package commands

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) blpop(args, commandAndArgs []string) resp.Value {
	if len(args) < 2 {
		return resp.SimpleError{Value: "BLPOP command must have at least 2 args"}
	}

	key := args[0]
	if _, ok := c.storage.StringStorage.Get(key); ok {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}
	timeoutMS, err := strconv.Atoi(args[1])
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("BLPOP command timeout (MS) argument atoi error: %v", err)}
	}
	poppedValue := c.storage.ListStorage.Blpop(key, timeoutMS)

	c.propagateWriteCommand(commandAndArgs)

	if poppedValue == nil {
		return resp.BulkString{Value: nil}
	}
	return resp.CreateBulkStringArray(key, *poppedValue)
}
