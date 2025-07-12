package state

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
)

type MasterServer struct {
	*BaseServer
	Replicas map[string]net.Conn
}

func NewMasterServer(args *config.Args, listener net.Listener) *MasterServer {
	ms := &MasterServer{
		BaseServer: NewBaseServer(args, listener),
		Replicas:   make(map[string]net.Conn),
	}
	ms.ReplicationInfo = replication.NewMasterInfo()
	ms.CommandController = commands.NewController(ms.Storage, ms.Args, ms.ReplicationInfo)
	return ms
}

func (ms *MasterServer) Start() {
	fmt.Println("START MASTER SERVER")
	ms.initStorage()
	ms.acceptConnections()
	ms.startExpiredKeysCleanup()
}

func (ms *MasterServer) IsMaster() bool {
	return true
}
