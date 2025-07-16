package servers

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
)

type master struct {
	*base
	replicas map[string]net.Conn
}

func newMaster(args *config.Args) *master {
	m := &master{
		base:     newBase(args),
		replicas: make(map[string]net.Conn),
	}
	m.CommandController = commands.NewController(m.Storage, m.Args, m)
	return m
}

func (m *master) Start() {
	fmt.Println("START MASTER SERVER")
	m.initStorage()
	m.acceptClientConnections()
	m.startExpiredKeysCleanup()
}

func (m *master) AddReplicaConn(addr string, replicaConn net.Conn) {
	m.replicas[addr] = replicaConn
}

func (m *master) Info() *replication.Info {
	return &replication.Info{
		Role:             "master",
		MasterReplID:     replication.GenerateReplicationId(),
		MasterReplOffset: 0,
	}
}
