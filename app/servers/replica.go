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

var _ replication.Replica = (*replica)(nil)

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

func (r *replica) ReadValueFromMaster() (resp.Value, error) {
	b := make([]byte, 1024)
	n, err := r.masterConn.Read(b)
	if err != nil {
		return nil, err
	}
	//Reading only first acceptable command from buffer, others are discarded
	_, value, err := r.RESPController.Decode(b[:n])
	return value, err
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
	r.processMasterHandshakePING()
	r.processMasterHandshakeREPLCONF()
}

func (r *replica) processMasterHandshakePING() {
	pingCommand := resp.CreateArray("PING")
	err := r.CommandController.Write(pingCommand, r.masterConn)
	if err != nil {
		log.Fatalf("Master handshake PING (1/3) write error: %s\n", err)
	}
	pingResult, err := r.ReadValueFromMaster()
	if err != nil {
		log.Fatalf("Master handshake PING (1/3) read value from master error: %s\n", err)
	}
	resp.AssertEqualSimpleString(pingResult, "PONG")
}

func (r *replica) processMasterHandshakeREPLCONF() {
	commands := []resp.Array{
		resp.CreateArray("REPLCONF", "listening-port", strconv.Itoa(r.Args.Port)),
		resp.CreateArray("REPLCONF", "capa", "psync2"),
	}

	for _, command := range commands {
		err := r.CommandController.Write(command, r.masterConn)
		if err != nil {
			log.Fatalf("Master handshake REPLCONF (2/3) write error: %s\n", err)
		}
		replconfResult, err := r.ReadValueFromMaster()
		if err != nil {
			log.Fatalf("Master handshake REPLCONF (2/3) read value from master error: %s\n", err)
		}
		resp.AssertEqualSimpleString(replconfResult, "OK")
	}
}
