package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func extractCommandAndArgs(commandAndArgs []resp.Value) ([]string, error) {
	result := make([]string, 0, len(commandAndArgs))

	for i, unit := range commandAndArgs {
		switch u := unit.(type) {
		case resp.BulkString:
			if u.Value == nil {
				return nil, fmt.Errorf("element %d is a null bulk string", i)
			}
			result = append(result, *u.Value)
		case resp.SimpleString:
			result = append(result, u.Value)
		case resp.Integer:
			result = append(result, strconv.Itoa(u.Value))
		default:
			return nil, fmt.Errorf("element %d is not a RESP bulk string or simple string or integer, got %T", i, unit)
		}
	}

	return result, nil
}

func isMastersBacklogBufferCommand(cmd resp.Value) bool {
	switch v := cmd.(type) {
	case resp.Array:
		if len(v.Value) == 0 {
			return false
		}
		cmdName, ok := v.Value[0].(resp.BulkString)
		if !ok {
			return false
		}
		uppercasedCmdName := strings.ToUpper(*cmdName.Value)
		return isPropagatedCommand(uppercasedCmdName) || isSpecialCommand(uppercasedCmdName)
	}
	return false
}

func isPropagatedCommand(cmd string) bool {
	return cmd == "SET"
}

func isSpecialCommand(cmd string) bool {
	return cmd == "PING"
}
