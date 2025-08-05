package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) discard(args []string, conn net.Conn) resp.Value {
	if len(args) != 0 {
		return resp.SimpleError{Value: "DISCARD command doesn't have args"}
	}
	if !c.transactionController.InTransaction(conn) {
		return resp.SimpleError{Value: "ERR DISCARD without MULTI"}
	}

	c.transactionController.RemoveConn(conn)
	return resp.SimpleString{Value: "OK"}
}
