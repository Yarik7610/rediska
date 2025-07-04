package commands

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func Echo(args []string) resp.BulkString {
	res := strings.Join(args, " ")
	return resp.BulkString{Value: &res}
}
