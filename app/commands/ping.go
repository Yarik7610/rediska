package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) ping(conn net.Conn) resp.Value {
	subscribeModePong := "pong"
	subscribeModeEmptyStr := ""

	if c.pubsubController.InSubscribeMode(conn) {
		return resp.Array{Value: []resp.Value{resp.BulkString{Value: &subscribeModePong}, resp.BulkString{Value: &subscribeModeEmptyStr}}}
	}
	return resp.SimpleString{Value: "PONG"}
}
