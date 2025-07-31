package commands

import (
	"fmt"
	"strconv"

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

func extractKeyValuesToStringSlice(m map[string]string) []string {
	keyValues := make([]string, 0)
	for key, value := range m {
		keyValues = append(keyValues, key, value)
	}
	return keyValues
}
