package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) publish(args []string, conn net.Conn) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "PUBLISH command must have 2 args"}
	}

	channelSubscribersCount := c.pubsubController.Publish(args[0], args[1])
	return resp.Integer{Value: channelSubscribersCount}
}
