package commands

import (
	"fmt"
	"log"
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

	results := make([]resp.Value, 0)
	commands, err := c.transactionController.GetQueue(conn)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("ERR %s", err)}
	}

	c.transactionController.RemoveConn(conn)

	for _, command := range commands {
		result, err := c.HandleCommand(command, conn, false)
		if err != nil {
			log.Printf("handle command error: %v, continue to work", err)
		}
		results = append(results, result)
	}

	return resp.Array{Value: results}
}
