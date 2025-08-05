package transaction

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

type Controller interface {
	AddConn(conn net.Conn)
	RemoveConn(conn net.Conn)
	InTransaction(conn net.Conn) bool
	EnqueueCommand(conn net.Conn, cmd resp.Value) error
	DequeueCommand(conn net.Conn) (resp.Value, error)
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

func (tc *controller) AddConn(conn net.Conn) {
	addr := utils.GetRemoteAddr(conn)
	if _, ok := tc.connQueues[addr]; ok {
		return
	}
	tc.connQueues[addr] = make([]resp.Value, 0)
}

func (tc *controller) RemoveConn(conn net.Conn) {
	addr := utils.GetRemoteAddr(conn)
	if _, ok := tc.connQueues[addr]; !ok {
		return
	}
	delete(tc.connQueues, addr)
}

func (tc *controller) InTransaction(conn net.Conn) bool {
	addr := utils.GetRemoteAddr(conn)
	_, ok := tc.connQueues[addr]
	return ok
}

func (tc *controller) EnqueueCommand(conn net.Conn, cmd resp.Value) error {
	addr := utils.GetRemoteAddr(conn)
	queue := tc.getCommandsQueue(conn)
	if queue == nil {
		return fmt.Errorf("conn %s isn't in transaction", addr)
	}
	tc.connQueues[addr] = append(queue, cmd)
	return nil
}

func (tc *controller) DequeueCommand(conn net.Conn) (resp.Value, error) {
	addr := utils.GetRemoteAddr(conn)
	queue := tc.getCommandsQueue(conn)
	if queue == nil {
		return nil, fmt.Errorf("conn %s isn't in transaction", addr)
	}
	if len(queue) == 0 {
		return nil, fmt.Errorf("conn %s, pop from empty queue detected", addr)
	}
	cmd := queue[0]
	tc.connQueues[addr] = queue[1:]
	return cmd, nil
}
