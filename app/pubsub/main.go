package pubsub

import (
	"log"
	"net"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

type subscriber struct {
	conn         net.Conn
	subscribedTo map[string]bool
}

type Controller interface {
	Publish(channel, message string) int
	Subscribe(conn net.Conn, channels ...string) []ChanAndLen
	Unsubscribe(conn net.Conn, channels ...string) []ChanAndLen
	InSubscribeMode(conn net.Conn) bool
	ValidateSubscribeModeCommand(cmd string, conn net.Conn) error
	UnsubscribeFromAllChannels(conn net.Conn)
}

type controller struct {
	channelSubs map[string][]*subscriber
	connSubs    map[string]*subscriber
	rwMut       sync.RWMutex
}

func NewController() Controller {
	return &controller{
		channelSubs: make(map[string][]*subscriber),
		connSubs:    make(map[string]*subscriber),
	}
}

type ChanAndLen struct {
	Channel         string
	SubscribedToLen int
}

func (c *controller) Publish(channel, message string) int {
	c.rwMut.Lock()
	defer c.rwMut.Unlock()

	channelSubs, ok := c.channelSubs[channel]
	if !ok {
		return 0
	}

	for _, sub := range channelSubs {
		go func(sub *subscriber) {
			err := writeMessageToSubscriber(channel, message, sub)
			if err != nil {
				log.Printf("Publishing error: %s", err)
			}
		}(sub)
	}

	return len(channelSubs)
}

func (c *controller) Subscribe(conn net.Conn, channels ...string) []ChanAndLen {
	sub := c.getOrCreateSubscriber(conn)

	c.rwMut.Lock()
	defer c.rwMut.Unlock()

	response := make([]ChanAndLen, 0)
	for _, channel := range channels {
		if !sub.subscribedTo[channel] {
			sub.subscribedTo[channel] = true
			c.channelSubs[channel] = append(c.channelSubs[channel], sub)
		}
		response = append(response, ChanAndLen{
			Channel:         channel,
			SubscribedToLen: len(sub.subscribedTo),
		})
	}
	return response
}

func (c *controller) Unsubscribe(conn net.Conn, channels ...string) []ChanAndLen {
	sub := c.getOrCreateSubscriber(conn)

	c.rwMut.Lock()
	defer c.rwMut.Unlock()

	response := make([]ChanAndLen, 0)
	for _, channel := range channels {
		c.channelSubs[channel] = c.removeSubscriberFromChannel(sub, channel)
		if len(c.channelSubs[channel]) == 0 {
			delete(c.channelSubs, channel)
		}

		response = append(response, ChanAndLen{
			Channel:         channel,
			SubscribedToLen: len(sub.subscribedTo),
		})
	}
	return response
}

func (c *controller) InSubscribeMode(conn net.Conn) bool {
	s := c.getSubscriber(conn)
	return s != nil
}

func (c *controller) UnsubscribeFromAllChannels(conn net.Conn) {
	subscriber := c.getSubscriber(conn)
	if subscriber == nil {
		return
	}

	c.rwMut.Lock()
	defer c.rwMut.Unlock()

	subscriberAddr := utils.GetRemoteAddr(subscriber.conn)

	for channel := range subscriber.subscribedTo {
		c.channelSubs[channel] = c.removeSubscriberFromChannel(subscriber, channel)
		if len(c.channelSubs[channel]) == 0 {
			delete(c.channelSubs, channel)
		}
	}

	delete(c.connSubs, subscriberAddr)
}
