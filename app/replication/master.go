package replication

import (
	"net"
)

type Master interface {
	Base
	AddReplicaConn(addr string, replicaConn net.Conn)
	GetReplicas() map[string]net.Conn
	SendRDBFile(replicaConn net.Conn)
	Propagate(args []string)
}
