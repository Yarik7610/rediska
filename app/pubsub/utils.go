package pubsub

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

func (subs *subscribers) removeSubscriberFromChannelSubs(sub *subscriber, channel string) []*subscriber {
	subsChan := subs.channelSubs[channel]
	result := make([]*subscriber, 0, len(subsChan))
	for _, s := range subsChan {
		if s != sub {
			result = append(result, s)
		}
	}
	return result
}

func (subs *subscribers) getSubscriber(conn net.Conn) *subscriber {
	subs.rwMut.RLock()
	defer subs.rwMut.RUnlock()

	addr := utils.GetRemoteAddr(conn)
	if s, ok := subs.connSubs[addr]; ok {
		return s
	}
	return nil
}

func (subs *subscribers) getOrCreateSubscriber(conn net.Conn) *subscriber {
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

	s := &subscriber{
		conn:         conn,
		ch:           make(chan string),
		subscribedTo: make(map[string]bool),
	}
	subs.connSubs[addr] = s
	return s
}
