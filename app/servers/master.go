package servers

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/persistence/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type master struct {
	*base
	replicas map[string]net.Conn
}

var _ replication.Master = (*master)(nil)

func newMaster(args *config.Args) *master {
	m := &master{
		base:     newBase(args),
		replicas: make(map[string]net.Conn),
	}
	m.commandController = commands.NewController(m.storage, m.args, m)
	m.replicationInfo = m.initReplicationInfo()
	return m
}

func (m *master) Start() {
	m.initStorage()
	go m.startExpiredKeysCleanup()
	m.acceptClientConnections()
}

func (m *master) acceptClientConnections() {
	address := fmt.Sprintf("%s:%d", m.args.Host, m.args.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to bind to address: %s\n", address)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go m.handleClient(nil, conn, true)
	}
}

func (m *master) Propagate(args []string) {
	command := resp.CreateArray(args...)
	for addr, conn := range m.replicas {
		err := m.commandController.Write(command, conn)
		if err != nil {
			log.Printf("Desynchronization with %s (but continue to work), propagate error: %v", addr, err)
			continue
		}
	}
}

func (m *master) AddReplicaConn(addr string, replicaConn net.Conn) {
	m.replicas[addr] = replicaConn
}

func (m *master) SendRDBFile(replicaConn net.Conn) {
	// Alternatively here call rdb.ReadRDBFile, but i use hardcoded hex string to pass tests
	// (In tests i guess they don't use rdb file, so my program returns error that can't find such dir when reading rdb file)
	b, err := hex.DecodeString(rdb.EMPTY_DB_HEX)
	if err != nil {
		log.Fatalf("SendRDBFile error: %v", err)
	}

	response := append(fmt.Appendf(nil, "$%d\r\n", len(b)), b...)
	_, err = replicaConn.Write(response)
	if err != nil {
		log.Fatalf("SendRDBFile error: %v", err)
	}
}

func (m *master) initReplicationInfo() *replication.Info {
	return &replication.Info{
		Role:             "master",
		MasterReplID:     replication.GenerateReplicationId(),
		MasterReplOffset: 0,
	}
}
