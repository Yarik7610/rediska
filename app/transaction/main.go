package transaction

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

type Controller interface {
	InTransaction(conn net.Conn) bool
	EnqueueCommand(conn net.Conn, cmd resp.Value)
	DequeueCommand(conn net.Conn) (resp.Value, error)
	RemoveConn(conn net.Conn)
}

// Don't use mutex for queues because, clients never cross and share data
// Mutex will block multiple clients for appeding commands to their own tx
type controller struct {
	connQueues map[string][]resp.Value
}

func NewController() Controller {
	return &controller{
		connQueues: make(map[string][]resp.Value),
	}
}

func (tc *controller) InTransaction(conn net.Conn) bool {
	addr := utils.GetRemoteAddr(conn)
	_, ok := tc.connQueues[addr]
	return ok
}

func (tc *controller) EnqueueCommand(conn net.Conn, cmd resp.Value) {
	addr := utils.GetRemoteAddr(conn)
	queue := tc.getOrCreateCommandsQueue(conn)
	tc.connQueues[addr] = append(queue, cmd)
}

func (tc *controller) DequeueCommand(conn net.Conn) (resp.Value, error) {
	queue := tc.getOrCreateCommandsQueue(conn)
	if len(queue) == 0 {
		return nil, fmt.Errorf("pop from empty queue detected")
	}
	addr := utils.GetRemoteAddr(conn)
	cmd := queue[0]
	tc.connQueues[addr] = queue[1:]
	return cmd, nil
}

func (tc *controller) RemoveConn(conn net.Conn) {
	addr := utils.GetRemoteAddr(conn)
	if _, ok := tc.connQueues[addr]; !ok {
		return
	}
	delete(tc.connQueues, addr)
}
