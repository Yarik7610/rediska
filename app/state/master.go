package state

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
)

type MasterServer struct {
	*BaseServer
	Replicas map[string]net.Conn
}

func NewMasterServer(args *config.Args) *MasterServer {
	ms := &MasterServer{
		BaseServer: NewBaseServer(args),
		Replicas:   make(map[string]net.Conn),
	}
	ms.CommandController = commands.NewController(ms.BaseServer.Storage, ms.BaseServer.Args, ms.IsMaster())
	return ms
}

func (ms *MasterServer) Start() {
	fmt.Println("START MASTER SERVER")
	ms.initStorage()
	ms.startExpiredKeysCleanup()
}

func (ms *MasterServer) IsMaster() bool {
	return true
}
