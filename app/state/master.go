package state

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
)

type masterServer struct {
	*baseServer
	Replicas map[string]net.Conn
}

func newMasterServer(args *config.Args, listener net.Listener) *masterServer {
	ms := &masterServer{
		baseServer: newBaseServer(args, listener),
		Replicas:   make(map[string]net.Conn),
	}
	ms.ReplicationInfo = newMasterInfo()
	ms.CommandController = commands.NewController(ms.Storage, ms.Args, ms.ReplicationInfo)
	return ms
}

func (ms *masterServer) Start() {
	fmt.Println("START MASTER SERVER")
	ms.initStorage()
	ms.acceptConnections()
	ms.startExpiredKeysCleanup()
}

func newMasterInfo() *replication.Info {
	return &replication.Info{
		Role:             "master",
		MasterReplID:     replication.GenerateReplicationId(),
		MasterReplOffset: 0,
	}
}
