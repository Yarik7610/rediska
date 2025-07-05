package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func extractCommandAndArgs(arr []resp.Value) ([]string, error) {
	result := make([]string, 0, len(arr))

	for i, unit := range arr {
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

func getExpiry(expireMark, expireValue string) (time.Duration, error) {
	atoiExpireValue, err := strconv.Atoi(expireValue)
	if err != nil {
		return 0, fmt.Errorf("expire value atoi error: %v", err)
	}

	expireDurationValue := time.Duration(atoiExpireValue)

	upperCasedExpireMark := strings.ToUpper(expireMark)
	switch upperCasedExpireMark {
	case "EX":
		return time.Second * expireDurationValue, nil
	case "PX":
		return time.Millisecond * expireDurationValue, nil
	default:
		return 0, fmt.Errorf("unknown expire mark: %s", upperCasedExpireMark)
	}
}
