package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func Ping() resp.SimpleString {
	return resp.SimpleString{Value: "PONG"}
}
