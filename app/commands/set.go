package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/state"
)

func Set(args []string, server *state.Server) resp.Value {
	if len(args) < 2 {
		return resp.SimpleError{Value: "SET command must have at least 2 args"}
	}

	key := args[0]
	value := args[1]

	if len(args) > 2 {
		expireMark := args[2]
		expireValue := args[3]
		expiry, err := getExpiry(expireMark, expireValue)
		if err != nil {
			return resp.SimpleError{Value: fmt.Sprintf("SET command get expiry error: %v", err)}
		}
		server.Storage.SetWithExpiry(key, value, expiry)
		return resp.SimpleString{Value: "OK"}
	}

	server.Storage.Set(key, value)
	return resp.SimpleString{Value: "OK"}
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
