package commands

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/pubsub"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) subcribeOrUnsubscribe(commandAndArgs []string, conn net.Conn) resp.Value {
	commandName := commandAndArgs[0]
	args := commandAndArgs[1:]
	if len(args) < 1 {
		return resp.SimpleError{Value: fmt.Sprintf("%s command must have at least 1 arg", strings.ToUpper(commandName))}
	}

	var gotResponses []pubsub.ChanAndLen
	if commandName == "SUBSCRIBE" {
		gotResponses = c.pubsubController.Subscribe(conn, args...)
	} else {
		gotResponses = c.pubsubController.Unsubscribe(conn, args...)
	}

	if len(gotResponses) == 1 {
		return pubsub.CreateRESPChannelAndLenResponse(strings.ToLower(commandName), gotResponses[0])
	}

	multipleRESPResponses := make([]resp.Value, 0)
	for _, gotResponse := range gotResponses {
		multipleRESPResponses = append(multipleRESPResponses, pubsub.CreateRESPChannelAndLenResponse(strings.ToLower(commandName), gotResponse))
	}
	return resp.Array{Value: multipleRESPResponses}
}

func (c *controller) publish(args []string) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "PUBLISH command must have 2 args"}
	}

	channelSubscribersCount := c.pubsubController.Publish(args[0], args[1])
	return resp.Integer{Value: channelSubscribersCount}
}
