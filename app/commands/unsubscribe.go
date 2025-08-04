package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/pubsub"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) unsubscribe(args []string, conn net.Conn) resp.Value {
	if len(args) < 1 {
		return resp.SimpleError{Value: "UNSUBSCRIBE command must have at least 1 arg"}
	}

	gotResponses := c.pubsubController.Unsubscribe(conn, args...)

	if len(gotResponses) == 1 {
		return pubsub.CreateRESPChannelAndLenResponse("unsubscribe", gotResponses[0])
	}

	multipleRESPSubs := make([]resp.Value, 0)
	for _, gotResponse := range gotResponses {
		multipleRESPSubs = append(multipleRESPSubs, pubsub.CreateRESPChannelAndLenResponse("unsubscribe", gotResponse))
	}
	return resp.Array{Value: multipleRESPSubs}
}
