package replication

import (
	"net"
)

type Master interface {
	Base
	AddReplicaConn(addr string, replicaConn net.Conn)
	GetReplicas() map[string]net.Conn
	IsReplica(conn net.Conn) bool
	SendRDBFile(replicaConn net.Conn)
	Propagate(args []string)
	GetAckCh() chan Ack
	SendAck(addr string, offset int)
	SetHasPendingWrites(val bool)
	HasPendingWrites() bool
}
