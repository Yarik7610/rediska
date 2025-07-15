package servers

import (
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
)

type Master struct {
	*Base
	replicas map[string]net.Conn
}

func newMaster(args *config.Args) *Master {
	m := &Master{
		Base:     newBase(args),
		replicas: make(map[string]net.Conn),
	}
	m.CommandController = commands.NewController(m.Storage, m.Args, m)
	return m
}

func (m *Master) Start() {
	fmt.Println("START MASTER SERVER")
	m.initStorage()
	m.acceptClientConnections()
	m.startExpiredKeysCleanup()
}

func (m *Master) AddReplicaConn(addr string, replicaConn net.Conn) {
	log.Println("ADDR TO ADD TO TABLE", addr)
	m.replicas[addr] = replicaConn
}

func (m *Master) Info() *replication.Info {
	return &replication.Info{
		Role:             "master",
		MasterReplID:     replication.GenerateReplicationId(),
		MasterReplOffset: 0,
	}
}
