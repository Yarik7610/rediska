package pubsub

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

func CreateRESPChannelAndLenResponse(action string, chanAndLen ChanAndLen) resp.Array {
	return resp.Array{Value: []resp.Value{
		resp.BulkString{Value: &action},
		resp.BulkString{Value: &chanAndLen.Channel},
		resp.Integer{Value: chanAndLen.SubscribedToLen},
	}}
}

func writeMessageToSubscriber(channel, message string, sub *subscriber) error {
	addr := utils.GetRemoteAddr(sub.conn)

	response := resp.CreateBulkStringArray("message", channel, message)
	b, err := response.Encode()
	if err != nil {
		return err
	}

	_, err = sub.conn.Write(b)
	if err != nil {
		return fmt.Errorf("write to subscriber %s error: %s", addr, err)
	}
	return nil
}

func (c *controller) removeSubscriberFromChannel(sub *subscriber, channel string) []*subscriber {
	delete(sub.subscribedTo, channel)
	subsChan := c.channelSubs[channel]
	result := make([]*subscriber, 0, len(subsChan))
	for _, s := range subsChan {
		if s != sub {
			result = append(result, s)
		}
	}
	return result
}

func (c *controller) getSubscriber(conn net.Conn) *subscriber {
	c.rwMut.RLock()
	defer c.rwMut.RUnlock()

	addr := utils.GetRemoteAddr(conn)
	if s, ok := c.connSubs[addr]; ok {
		return s
	}
	return nil
}

func (c *controller) getOrCreateSubscriber(conn net.Conn) *subscriber {
	c.rwMut.RLock()
	addr := utils.GetRemoteAddr(conn)

	if s, ok := c.connSubs[addr]; ok {
		c.rwMut.RUnlock()
		return s
	}
	c.rwMut.RUnlock()

	c.rwMut.Lock()
	defer c.rwMut.Unlock()

	if s, ok := c.connSubs[addr]; ok {
		return s
	}

	s := &subscriber{
		conn:         conn,
		subscribedTo: make(map[string]bool),
	}
	c.connSubs[addr] = s
	return s
}
