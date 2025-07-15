package state

import (
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type replicaServer struct {
	*baseServer
	masterConn net.Conn
}

func newReplicaServer(args *config.Args) *replicaServer {
	rs := &replicaServer{
		baseServer: newBaseServer(args),
	}
	rs.ReplicationInfo = newReplicaInfo()
	rs.CommandController = commands.NewController(rs.Storage, rs.Args, rs.ReplicationInfo)
	return rs
}

func (rs *replicaServer) Start() {
	fmt.Println("START REPLICA SERVER")
	rs.initStorage()
	rs.connectToMaster()
	rs.acceptClientConnections()
}

func newReplicaInfo() *replication.Info {
	return &replication.Info{
		Role:             "slave",
		MasterReplID:     replication.GenerateReplicationId(),
		MasterReplOffset: 0,
	}
}

func (rs *replicaServer) connectToMaster() {
	rs.dialMaster()
	rs.processMasterHandshake()
}

func (rs *replicaServer) dialMaster() {
	address := fmt.Sprintf("%s:%d", rs.Args.ReplicaOf.Host, rs.Args.ReplicaOf.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to dial master address: %s\n", address)
	}
	rs.masterConn = conn
}

func (rs *replicaServer) processMasterHandshake() {
	pingCommand := resp.Array{Value: []resp.Value{resp.BulkString{Value: resp.StrPtr("PING")}}}
	err := rs.CommandController.EncodeAndWrite(pingCommand, rs.masterConn)
	if err != nil {
		log.Fatalf("Master handshake PING (1/3) error: %s\n", err)
	}
}
