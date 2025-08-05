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
	IsTransactionCommand(cmd string) bool
	EnqueueCommand(conn net.Conn, cmd resp.Value) error
	GetQueue(conn net.Conn) ([]resp.Value, error)
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

func (c *controller) AddConn(conn net.Conn) {
	addr := utils.GetRemoteAddr(conn)
	if _, ok := c.connQueues[addr]; ok {
		return
	}
	c.connQueues[addr] = make([]resp.Value, 0)
}

func (c *controller) RemoveConn(conn net.Conn) {
	addr := utils.GetRemoteAddr(conn)
	if _, ok := c.connQueues[addr]; !ok {
		return
	}
	delete(c.connQueues, addr)
}

func (c *controller) InTransaction(conn net.Conn) bool {
	addr := utils.GetRemoteAddr(conn)
	_, ok := c.connQueues[addr]
	return ok
}

func (c *controller) EnqueueCommand(conn net.Conn, cmd resp.Value) error {
	addr := utils.GetRemoteAddr(conn)
	queue := c.getCommandsQueue(conn)
	if queue == nil {
		return fmt.Errorf("conn %s isn't in transaction", addr)
	}
	c.connQueues[addr] = append(queue, cmd)
	return nil
}

func (c *controller) GetQueue(conn net.Conn) ([]resp.Value, error) {
	addr := utils.GetRemoteAddr(conn)
	queue := c.getCommandsQueue(conn)
	if queue == nil {
		return nil, fmt.Errorf("conn %s isn't in transaction", addr)
	}
	return queue, nil
}
