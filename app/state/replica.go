package state

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
)

type ReplicaServer struct {
	*BaseServer
}

func NewReplicaServer(args *config.Args, listener net.Listener) *ReplicaServer {
	rs := &ReplicaServer{
		BaseServer: NewBaseServer(args, listener),
	}
	rs.ReplicationInfo = replication.NewReplicaInfo()
	rs.CommandController = commands.NewController(rs.Storage, rs.Args, rs.ReplicationInfo)
	return rs
}

func (rs *ReplicaServer) Start() {
	fmt.Println("START REPLICA SERVER")
	rs.initStorage()
	rs.acceptConnections()
	rs.startExpiredKeysCleanup()
}

func (ms *ReplicaServer) IsMaster() bool {
	return false
}
