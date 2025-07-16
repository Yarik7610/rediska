package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (*Controller) ping() resp.SimpleString {
	return resp.SimpleString{Value: "PONG"}
}
