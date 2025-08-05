package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) multi(args []string, conn net.Conn) resp.Value {
	if len(args) != 0 {
		return resp.SimpleError{Value: "MULTI command doesn't have args"}
	}

	c.transactionController.AddConn(conn)
	return resp.SimpleString{Value: "OK"}
}
