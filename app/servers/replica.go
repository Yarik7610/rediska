package servers

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type replica struct {
	*base
	masterConn net.Conn
}

func newReplica(args *config.Args) *replica {
	r := &replica{
		base: newBase(args),
	}
	r.CommandController = commands.NewController(r.Storage, r.Args, r)
	return r
}

func (r *replica) Start() {
	fmt.Println("START REPLICA SERVER")
	r.initStorage()
	r.connectToMaster()
	r.acceptClientConnections()
}

func (*replica) Info() *replication.Info {
	return &replication.Info{
		Role:             "slave",
		MasterReplID:     replication.GenerateReplicationId(),
		MasterReplOffset: 0,
	}
}

func (r *replica) ReadFromMaster() ([]byte, error) {
	b := make([]byte, 1024)
	n, err := r.masterConn.Read(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}

func (r *replica) connectToMaster() {
	r.dialMaster()
	r.processMasterHandshake()
}

func (r *replica) dialMaster() {
	address := fmt.Sprintf("%s:%d", r.Args.ReplicaOf.Host, r.Args.ReplicaOf.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to dial master address: %s\n", address)
	}
	r.masterConn = conn
}

func (r *replica) processMasterHandshake() {
	pingCommand := resp.Array{Value: []resp.Value{resp.BulkString{Value: resp.StrPtr("PING")}}}
	err := r.CommandController.Write(pingCommand, r.masterConn)
	if err != nil {
		log.Fatalf("Master handshake PING (1/3) error: %s\n", err)
	}

	pingResult, err := r.ReadFromMaster()
	if err != nil {
		log.Fatalf("Master handshake PING (1/3) error: read from master error: %s\n", err)
	}
	_, value, err := r.RESPController.Decode(pingResult)
	if err != nil {
		log.Fatalf("Master handshake PING (1/3) error: decode error: %s", err)
	}
	v, ok := value.(resp.SimpleString)
	if !ok || v.Value != "PONG" {
		log.Fatalf("Master handshake PING (1/3) error: expected PONG, got: %v", value)
	}

	replconfCommand := resp.Array{Value: []resp.Value{
		resp.BulkString{Value: resp.StrPtr("REPLCONF")},
		resp.BulkString{Value: resp.StrPtr("listening-port")},
		resp.BulkString{Value: resp.StrPtr(strconv.Itoa(r.Args.Port))},
	}}
	err = r.CommandController.Write(replconfCommand, r.masterConn)
	if err != nil {
		log.Fatalf("Master handshake REPLCONF (2/3) error: %s\n", err)
	}
	replconfCommand = resp.Array{Value: []resp.Value{
		resp.BulkString{Value: resp.StrPtr("REPLCONF")},
		resp.BulkString{Value: resp.StrPtr("capa")},
		resp.BulkString{Value: resp.StrPtr("psync2")},
	}}
	err = r.CommandController.Write(replconfCommand, r.masterConn)
	if err != nil {
		log.Fatalf("Master handshake REPLCONF (2/3) error: %s\n", err)
	}
}
