package state

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
)

type ReplicaServer struct {
	*BaseServer
}

func NewReplicaServer(args *config.Args) *ReplicaServer {
	rs := &ReplicaServer{
		BaseServer: NewBaseServer(args),
	}
	rs.CommandController = commands.NewController(rs.BaseServer.Storage, rs.BaseServer.Args, rs.IsMaster())
	return rs
}

func (rs *ReplicaServer) Start() {
	fmt.Println("START REPLICA SERVER")
}

func (ms *ReplicaServer) IsMaster() bool {
	return false
}
