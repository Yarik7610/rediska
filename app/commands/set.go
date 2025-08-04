package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) set(args, commandAndArgs []string) resp.Value {
	if len(args) < 2 {
		return resp.SimpleError{Value: "SET command must have at least 2 args"}
	}

	key := args[0]
	value := args[1]

	if c.storage.KeyExistsWithOtherType(key, memory.TYPE_STRING) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	if len(args) > 2 {
		expireMark := args[2]
		expireValue := args[3]
		expiry, err := getExpireTime(expireMark, expireValue)
		if err != nil {
			return resp.SimpleError{Value: fmt.Sprintf("SET command get expiry error: %v", err)}
		}
		c.storage.StringStorage().SetWithExpiry(key, value, expiry)
		c.propagateWriteCommand(commandAndArgs)
		return resp.SimpleString{Value: "OK"}
	}

	c.storage.StringStorage().Set(key, value)
	c.propagateWriteCommand(commandAndArgs)
	return resp.SimpleString{Value: "OK"}
}

func getExpireTime(expireMark, expireDuration string) (time.Time, error) {
	atoiExpireDuration, err := strconv.Atoi(expireDuration)
	if err != nil {
		return time.Time{}, fmt.Errorf("expire duration atoi error: %v", err)
	}

	expireDurationValue := time.Duration(atoiExpireDuration)

	upperCasedExpireMark := strings.ToUpper(expireMark)
	switch upperCasedExpireMark {
	case "EX":
		return time.Now().Add(time.Second * expireDurationValue), nil
	case "PX":
		return time.Now().Add(time.Millisecond * expireDurationValue), nil
	default:
		return time.Time{}, fmt.Errorf("unknown expire mark: %s", upperCasedExpireMark)
	}
}
