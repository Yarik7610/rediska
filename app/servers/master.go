package servers

import (
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

type master struct {
	*base
	replicationController replication.MasterController
}

func newMaster(args *config.Args) Server {
	m := &master{
		base:                  newBase(args),
		replicationController: replication.NewMasterController(),
	}
	m.commandController = commands.NewController(m.args, m.storage, m.replicationController, m.pubsubController, m.transactionController)
	return m
}

func (m *master) Start() {
	m.initStorage()
	listener := m.listenTCP()
	go m.startExpiredStringKeysCleanup()
	m.acceptClientConnections(listener)
}

func (m *master) acceptClientConnections(listener net.Listener) {
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		m.handleClientWithCleanup(nil, conn, true)
	}
}

func (m *master) handleClientWithCleanup(initialBuffer []byte, conn net.Conn, writeResponseToConn bool) {
	go func() {
		defer m.cleanUpConn(conn)
		m.handleClient(initialBuffer, conn, writeResponseToConn)
	}()
}

func (m *master) cleanUpConn(conn net.Conn) {
	addr := utils.GetRemoteAddr(conn)
	m.replicationController.RemoveReplicaConn(addr)
}
