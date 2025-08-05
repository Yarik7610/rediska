package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) exec(args []string, conn net.Conn) resp.Value {
	if len(args) != 0 {
		return resp.SimpleError{Value: "EXEC command doesn't have args"}
	}

	if !c.transactionController.InTransaction(conn) {
		return resp.SimpleError{Value: "ERR EXEC without MULTI"}
	}

	return resp.SimpleString{Value: "OK"}
}
