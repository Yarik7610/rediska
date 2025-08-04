package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/pubsub"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) subscribe(args []string, conn net.Conn) resp.Value {
	if len(args) < 1 {
		return resp.SimpleError{Value: "SUBSCRIBE command must have at least 1 arg"}
	}

	gotResponses := c.pubsubController.Subscribe(conn, args...)

	if len(gotResponses) == 1 {
		return pubsub.CreateRESPChannelAndLenResponse("subscribe", gotResponses[0])
	}

	multipleRESPSubs := make([]resp.Value, 0)
	for _, gotResponse := range gotResponses {
		multipleRESPSubs = append(multipleRESPSubs, pubsub.CreateRESPChannelAndLenResponse("subscribe", gotResponse))
	}
	return resp.Array{Value: multipleRESPSubs}
}
