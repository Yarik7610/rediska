package replication

import "net"

type Replica interface {
	Base
	GetMasterConn() net.Conn
}
