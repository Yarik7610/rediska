package servers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"

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
	r.commandController = commands.NewController(r.storage, r.args, r)
	r.replicationInfo = r.initReplicationInfo()
	return r
}

func (r *replica) Start() {
	fmt.Println("START REPLICA SERVER")
	r.initStorage()
	r.connectToMaster()
	r.acceptClientConnections()
}

func (*replica) initReplicationInfo() *replication.Info {
	return &replication.Info{
		Role:             "slave",
		MasterReplID:     "?",
		MasterReplOffset: -1,
	}
}

func (r *replica) ReadFromMaster() ([]byte, int, error) {
	b := make([]byte, 1024)
	n, err := r.masterConn.Read(b)
	if errors.Is(err, io.EOF) {
		return b, n, nil
	}

	return b, n, err
}

func (r *replica) readValueFromMaster() (resp.Value, error) {
	//Reading only first acceptable command from buffer, others are discarded
	b, n, err := r.ReadFromMaster()
	if err != nil {
		return nil, fmt.Errorf("read from master error: %v", err)
	}
	_, value, err := r.respController.Decode(b[:n])
	return value, err
}

func (r *replica) connectToMaster() {
	r.dialMaster()
	r.processMasterHandshake()
}

func (r *replica) dialMaster() {
	address := fmt.Sprintf("%s:%d", r.args.ReplicaOf.Host, r.args.ReplicaOf.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to dial master address: %s\n", address)
	}
	r.masterConn = conn
}

func (r *replica) processMasterHandshake() {
	r.processMasterHandshakePING()
	r.processMasterHandshakeREPLCONF()
	r.processMasterHandshakePSYNC()
}

func (r *replica) processMasterHandshakePING() {
	pingCommand := resp.CreateArray("PING")
	err := r.commandController.Write(pingCommand, r.masterConn)
	if err != nil {
		log.Fatalf("Master handshake PING (1/3) write error: %s\n", err)
	}
	pingResult, err := r.readValueFromMaster()
	if err != nil {
		log.Fatalf("Master handshake PING (1/3) read value from master error: %s\n", err)
	}
	resp.AssertEqualSimpleString(pingResult, "PONG")
}

func (r *replica) processMasterHandshakeREPLCONF() {
	commands := []resp.Array{
		resp.CreateArray("REPLCONF", "listening-port", strconv.Itoa(r.args.Port)),
		resp.CreateArray("REPLCONF", "capa", "psync2"),
	}

	for _, command := range commands {
		err := r.commandController.Write(command, r.masterConn)
		if err != nil {
			log.Fatalf("Master handshake REPLCONF (2/3) write error: %s\n", err)
		}
		replconfResult, err := r.readValueFromMaster()
		if err != nil {
			log.Fatalf("Master handshake REPLCONF (2/3) read value from master error: %s\n", err)
		}
		resp.AssertEqualSimpleString(replconfResult, "OK")
	}
}

func (r *replica) processMasterHandshakePSYNC() {
	psyncCommand := resp.CreateArray("PSYNC", "?", "-1")
	err := r.commandController.Write(psyncCommand, r.masterConn)
	if err != nil {
		log.Fatalf("Master handshake PSYNC (3/3) write error: %s\n", err)
	}
	psyncResult, err := r.readValueFromMaster()
	if err != nil {
		log.Fatalf("Master handshake PSYNC (3/3) read value from master error: %s\n", err)
	}

	simpleString, ok := psyncResult.(resp.SimpleString)
	if !ok {
		log.Fatalf("Master handshake PSYNC (3/3) psync response has wrong RESP type, expected: %T got %T", reflect.TypeOf(resp.SimpleString{}), simpleString)
	}
	splitted := strings.Split(simpleString.Value, " ")
	replID := splitted[1]
	replOffset := splitted[2]

	atoiReplOffset, err := strconv.Atoi(replOffset)
	if err != nil {
		log.Fatalf("Master handshake PSYNC (3/3) psync response has wrong replication offset: %v", err)
	}
	r.SetMasterReplID(replID)
	r.SetMasterReplOfffset(atoiReplOffset)
}
