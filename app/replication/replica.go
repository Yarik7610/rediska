package replication

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Replica interface {
	Base
	GetMasterConn() net.Conn
	UpdateMasterReplOffsetWithCmd(cmd resp.Value)
}
