package replication

import "net"

type Master interface {
	Main
	AddReplicaConn(addr string, replicaConn net.Conn)
}
