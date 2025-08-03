package pubsub

import (
	"net"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

type Subscriber struct {
	Conn         net.Conn
	Ch           chan string
	SubscribedTo map[string]bool
}

type Subscribers struct {
	channelSubs map[string][]*Subscriber
	connSubs    map[string]*Subscriber
	rwMut       sync.RWMutex
}

func NewSubscribers() *Subscribers {
	return &Subscribers{
		channelSubs: make(map[string][]*Subscriber),
		connSubs:    make(map[string]*Subscriber),
	}
}

type SubscribeResponse struct {
	Channel         string
	SubscribedToLen int
}

func (subs *Subscribers) Subscribe(conn net.Conn, channels ...string) []SubscribeResponse {
	subscriber := subs.getOrCreateSubscriber(conn)

	subs.rwMut.Lock()
	defer subs.rwMut.Unlock()

	response := make([]SubscribeResponse, 0)
	for _, channel := range channels {
		if !subscriber.SubscribedTo[channel] {
			subscriber.SubscribedTo[channel] = true
			subs.channelSubs[channel] = append(subs.channelSubs[channel], subscriber)
		}
		response = append(response, SubscribeResponse{
			Channel:         channel,
			SubscribedToLen: len(subscriber.SubscribedTo),
		})
	}
	return response
}

func (subs *Subscribers) UnsubscribeFromAllChannels(conn net.Conn) {
	subscriber := subs.getSubscriber(conn)
	if subscriber == nil {
		return
	}

	subs.rwMut.Lock()
	defer subs.rwMut.Unlock()

	subscriberAddr := utils.GetRemoteAddr(subscriber.Conn)

	for channel := range subscriber.SubscribedTo {
		subs.channelSubs[channel] = subs.removeSubscriberFromChannelSubs(subscriber, channel)
		delete(subscriber.SubscribedTo, channel)

		if len(subs.channelSubs[channel]) == 0 {
			delete(subs.channelSubs, channel)
		}
	}

	delete(subs.connSubs, subscriberAddr)
}

func (subs *Subscribers) removeSubscriberFromChannelSubs(subscriber *Subscriber, channel string) []*Subscriber {
	subsChan := subs.channelSubs[channel]
	result := make([]*Subscriber, 0, len(subsChan)-1)
	for _, s := range subsChan {
		if s != subscriber {
			result = append(result, s)
		}
	}
	return result
}

func (subs *Subscribers) getSubscriber(conn net.Conn) *Subscriber {
	subs.rwMut.RLock()
	defer subs.rwMut.RUnlock()

	addr := utils.GetRemoteAddr(conn)
	if s, ok := subs.connSubs[addr]; ok {
		return s
	}
	return nil
}

func (subs *Subscribers) getOrCreateSubscriber(conn net.Conn) *Subscriber {
	subs.rwMut.RLock()
	addr := utils.GetRemoteAddr(conn)

	if s, ok := subs.connSubs[addr]; ok {
		subs.rwMut.RUnlock()
		return s
	}
	subs.rwMut.RUnlock()

	subs.rwMut.Lock()
	defer subs.rwMut.Unlock()

	if s, ok := subs.connSubs[addr]; ok {
		return s
	}

	s := &Subscriber{
		Conn:         conn,
		Ch:           make(chan string),
		SubscribedTo: make(map[string]bool),
	}
	subs.connSubs[addr] = s
	return s
}
