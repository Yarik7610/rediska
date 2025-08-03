package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/pubsub"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) subscribe(args []string, conn net.Conn) resp.Value {
	if len(args) < 1 {
		return resp.SimpleError{Value: "SUBSCRIBE command must have at least 1 arg"}
	}

	gotResponses := c.subscribers.Subscribe(conn, args...)

	if len(gotResponses) == 1 {
		return createRESPSubscribeChannelResponse(gotResponses[0])
	}

	multipleRESPSubs := make([]resp.Value, 0)
	for _, gotResponse := range gotResponses {
		multipleRESPSubs = append(multipleRESPSubs, createRESPSubscribeChannelResponse(gotResponse))
	}
	return resp.Array{Value: multipleRESPSubs}
}

func createRESPSubscribeChannelResponse(subscribeResponse pubsub.SubscribeResponse) resp.Array {
	action := "subscribe"
	return resp.Array{Value: []resp.Value{
		resp.BulkString{Value: &action},
		resp.BulkString{Value: &subscribeResponse.Channel},
		resp.Integer{Value: subscribeResponse.SubscribedToLen},
	}}
}
