package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) push(commandAndArgs []string) resp.Value {
	commandName := strings.ToUpper(commandAndArgs[0])
	args := commandAndArgs[1:]
	if len(args) < 2 {
		return resp.SimpleError{Value: fmt.Sprintf("%s command must have at least 2 args", commandName)}
	}

	key := args[0]
	if c.storage.KeyExistsWithOtherType(key, memory.TYPE_LIST) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	values := args[1:]
	len := 0
	if commandName == "RPUSH" {
		len = c.storage.ListStorage().Rpush(key, values...)
	} else {
		len = c.storage.ListStorage().Lpush(key, values...)
	}
	c.propagateWriteCommand(commandAndArgs)
	return resp.Integer{Value: len}
}

func (c *controller) pop(commandAndArgs []string) resp.Value {
	commandName := strings.ToUpper(commandAndArgs[0])
	args := commandAndArgs[1:]
	if len(args) < 1 {
		return resp.SimpleError{Value: fmt.Sprintf("%s command must have at least 1 arg", commandName)}
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
			return resp.SimpleError{Value: fmt.Sprintf("%s command count argument atoi error: %v", commandName, err)}
		}
	}

	var poppedValues []string
	if commandName == "RPOP" {
		poppedValues = c.storage.ListStorage().Rpop(key, count)
	} else {
		poppedValues = c.storage.ListStorage().Lpop(key, count)
	}
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

func (c *controller) bpop(commandAndArgs []string) resp.Value {
	commandName := strings.ToUpper(commandAndArgs[0])
	args := commandAndArgs[1:]
	if len(args) < 2 {
		return resp.SimpleError{Value: fmt.Sprintf("%s command must have at least 2 args", commandName)}
	}

	key := args[0]
	if c.storage.KeyExistsWithOtherType(key, memory.TYPE_LIST) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	timeoutS, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("%s command timeout (S) argument parseFloat error: %v", commandName, err)}
	}

	var poppedValue *string
	if commandName == "BRPOP" {
		poppedValue = c.storage.ListStorage().Brpop(key, timeoutS)
	} else {
		poppedValue = c.storage.ListStorage().Blpop(key, timeoutS)
	}

	c.propagateWriteCommand(commandAndArgs)

	if poppedValue == nil {
		return resp.BulkString{Value: nil}
	}
	return resp.CreateBulkStringArray(key, *poppedValue)
}

func (c *controller) lrange(args []string) resp.Value {
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

func (c *controller) llen(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "LLEN command must have 1 arg"}
	}

	key := args[0]
	if c.storage.KeyExistsWithOtherType(key, memory.TYPE_LIST) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	len := c.storage.ListStorage().Llen(key)
	return resp.Integer{Value: len}
}
