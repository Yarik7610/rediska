package transaction

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

func (tc *controller) getCommandsQueue(conn net.Conn) []resp.Value {
	addr := utils.GetRemoteAddr(conn)
	if q, ok := tc.connQueues[addr]; ok {
		return q
	}
	return nil
}
