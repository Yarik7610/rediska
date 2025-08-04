package pubsub

import (
	"net"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

type subscriber struct {
	conn         net.Conn
	ch           chan string
	subscribedTo map[string]bool
}

type Subscribers interface {
	Subscribe(conn net.Conn, channels ...string) []SubscribeResponse
	InSubscribeMode(conn net.Conn) bool
	HandleSubscribeModeCommand(cmd string, conn net.Conn) error
	UnsubscribeFromAllChannels(conn net.Conn)
}

type subscribers struct {
	channelSubs map[string][]*subscriber
	connSubs    map[string]*subscriber
	rwMut       sync.RWMutex
}

func NewSubscribers() Subscribers {
	return &subscribers{
		channelSubs: make(map[string][]*subscriber),
		connSubs:    make(map[string]*subscriber),
	}
}

type SubscribeResponse struct {
	Channel         string
	SubscribedToLen int
}

func (subs *subscribers) Subscribe(conn net.Conn, channels ...string) []SubscribeResponse {
	subscriber := subs.getOrCreateSubscriber(conn)

	subs.rwMut.Lock()
	defer subs.rwMut.Unlock()

	response := make([]SubscribeResponse, 0)
	for _, channel := range channels {
		if !subscriber.subscribedTo[channel] {
			subscriber.subscribedTo[channel] = true
			subs.channelSubs[channel] = append(subs.channelSubs[channel], subscriber)
		}
		response = append(response, SubscribeResponse{
			Channel:         channel,
			SubscribedToLen: len(subscriber.subscribedTo),
		})
	}
	return response
}

func (subs *subscribers) InSubscribeMode(conn net.Conn) bool {
	s := subs.getSubscriber(conn)
	return s != nil
}

func (subs *subscribers) UnsubscribeFromAllChannels(conn net.Conn) {
	subscriber := subs.getSubscriber(conn)
	if subscriber == nil {
		return
	}

	subs.rwMut.Lock()
	defer subs.rwMut.Unlock()

	subscriberAddr := utils.GetRemoteAddr(subscriber.conn)

	for channel := range subscriber.subscribedTo {
		subs.channelSubs[channel] = subs.removeSubscriberFromChannelSubs(subscriber, channel)
		delete(subscriber.subscribedTo, channel)

		if len(subs.channelSubs[channel]) == 0 {
			delete(subs.channelSubs, channel)
		}
	}

	delete(subs.connSubs, subscriberAddr)
}
