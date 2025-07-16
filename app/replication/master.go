package replication

import "net"

type Master interface {
	Base
	AddReplicaConn(addr string, replicaConn net.Conn)
}
