package state

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
)

type replicaServer struct {
	*baseServer
}

func newReplicaServer(args *config.Args, listener net.Listener) *replicaServer {
	rs := &replicaServer{
		baseServer: newBaseServer(args, listener),
	}
	rs.ReplicationInfo = newReplicaInfo()
	rs.CommandController = commands.NewController(rs.Storage, rs.Args, rs.ReplicationInfo)
	return rs
}

func (rs *replicaServer) Start() {
	fmt.Println("START REPLICA SERVER")
	rs.initStorage()
	rs.acceptConnections()
	rs.startExpiredKeysCleanup()
}

func newReplicaInfo() *replication.Info {
	return &replication.Info{
		Role:             "slave",
		MasterReplID:     replication.GenerateReplicationId(),
		MasterReplOffset: 0,
	}
}
