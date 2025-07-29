package commands

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) lrange(args []string) resp.Value {
	if len(args) != 3 {
		return resp.SimpleError{Value: "LRANGE command must have 3 args"}
	}

	key := args[0]
	if c.storage.KeyExistsWithOtherType(key, memory.TYPE_LIST) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	startIdx := args[1]
	stopIdx := args[2]

	startIdxAtoi, err := strconv.Atoi(startIdx)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("LRANGE command start atoi error: %v", err)}
	}
	stopIdxAtoi, err := strconv.Atoi(stopIdx)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("LRANGE command stop atoi error: %v", err)}
	}

	values := c.storage.ListStorage().Lrange(key, startIdxAtoi, stopIdxAtoi)
	return resp.CreateBulkStringArray(values...)
}
