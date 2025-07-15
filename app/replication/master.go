package replication

import "net"

type Master interface {
	Replication
	AddReplicaConn(addr string, replicaConn net.Conn)
}
