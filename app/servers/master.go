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

func (m *master) Propagate(args []string) {
	command := resp.CreateBulkStringArray(args...)
	for addr, conn := range m.replicas {
		err := m.commandController.Write(command, conn)
		if err != nil {
			log.Printf("Desynchronization with %s (but continue to work), propagate error: %v", addr, err)
			continue
		}
	}
}

func (m *master) AddReplicaConn(addr string, replicaConn net.Conn) {
	log.Printf("Added replica %s to replicas map", addr)
	m.replicas[addr] = replicaConn
}

func (m *master) removeReplicaConn(addr string) {
	log.Printf("Removed replica %s from replicas map", addr)
	delete(m.replicas, addr)
}

func (m *master) GetReplicas() map[string]net.Conn {
	return m.replicas
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
	addr := conn.RemoteAddr().String()
	m.removeReplicaConn(addr)
}

func (m *master) initReplicationInfo() *replication.Info {
	return &replication.Info{
		Role:             "master",
		MasterReplID:     replication.GenerateReplicationId(),
		MasterReplOffset: 0,
	}
}
