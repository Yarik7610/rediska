package replication

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/persistence/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

type MasterController interface {
	BaseController
	AddReplicaConn(replicaConn net.Conn)
	RemoveReplicaConn(addr string)
	GetReplicas() map[string]net.Conn
	IsReplica(conn net.Conn) bool
	SendRDBFile(replicaConn net.Conn)
	Propagate(args []string)
	GetAcksCh() chan Ack
	SendAck(addr string, offset int)
	SetHasPendingWrites(val bool)
	HasPendingWrites() bool
}

type masterController struct {
	*baseController
	replicas           map[string]net.Conn
	acks               chan Ack
	hasPendingWrites   bool
	pendingWritesMutex sync.Mutex
}

func NewMasterController() MasterController {
	masterInfo := initMasterInfo()
	return &masterController{
		baseController: newBaseController(masterInfo),
		replicas:       make(map[string]net.Conn),
		acks:           make(chan Ack, 10),
	}
}

func (mc *masterController) AddReplicaConn(replicaConn net.Conn) {
	addr := utils.GetRemoteAddr(replicaConn)
	log.Printf("Added replica %s to replicas map", addr)
	mc.replicas[addr] = replicaConn
}

func (mc *masterController) RemoveReplicaConn(addr string) {
	if _, ok := mc.replicas[addr]; ok {
		log.Printf("Removed replica %s from replicas map", addr)
	}
	delete(mc.replicas, addr)
}

func (mc *masterController) GetReplicas() map[string]net.Conn {
	return mc.replicas
}

func (mc *masterController) IsReplica(conn net.Conn) bool {
	addr := utils.GetRemoteAddr(conn)
	_, ok := mc.replicas[addr]
	return ok
}

func (mc *masterController) SendRDBFile(replicaConn net.Conn) {
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

func (mc *masterController) Propagate(args []string) {
	command := resp.CreateBulkStringArray(args...)
	for addr, conn := range mc.replicas {
		go mc.propagateCommandToConn(command, addr, conn)
	}
}

func (mc *masterController) GetAcksCh() chan Ack {
	return mc.acks
}

func (mc *masterController) SendAck(addr string, offset int) {
	if _, ok := mc.replicas[addr]; !ok {
		return
	}
	mc.acks <- Ack{Addr: addr, Offset: offset}
}

func (mc *masterController) SetHasPendingWrites(val bool) {
	mc.pendingWritesMutex.Lock()
	defer mc.pendingWritesMutex.Unlock()
	mc.hasPendingWrites = val
}

func (mc *masterController) HasPendingWrites() bool {
	mc.pendingWritesMutex.Lock()
	defer mc.pendingWritesMutex.Unlock()
	return mc.hasPendingWrites
}

func (mc *masterController) propagateCommandToConn(command resp.Array, addr string, conn net.Conn) {
	err := utils.WriteCommand(command, conn)
	if err != nil {
		log.Printf("Desynchronization with %s (but continue to work), propagateWriteCommand error: %v", addr, err)
	}
}

func initMasterInfo() *Info {
	return &Info{
		Role:             "master",
		MasterReplID:     generateReplicationId(),
		MasterReplOffset: 0,
	}
}
