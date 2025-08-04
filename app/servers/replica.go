package servers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type replica struct {
	*base
	masterConn         net.Conn
	acceptClientsReady chan struct{}
	connectionsWG      sync.WaitGroup
	masterConnBuffer   []byte
}

var _ replication.Replica = (*replica)(nil)

func newReplica(args *config.Args) *replica {
	r := &replica{
		base:               newBase(args),
		acceptClientsReady: make(chan struct{}),
		masterConnBuffer:   make([]byte, 0),
	}
	r.commandController = commands.NewController(r.storage, r.args, r.pubsubController, r)
	r.replicationInfo = r.initReplicationInfo()
	return r
}

func (r *replica) Start() {
	r.initStorage()

	r.connectionsWG.Add(1)
	go func() {
		defer r.connectionsWG.Done()
		r.acceptClientConnections()
	}()

	<-r.acceptClientsReady

	r.connectionsWG.Add(1)
	go func() {
		defer r.connectionsWG.Done()
		r.connectToMaster()
	}()

	r.connectionsWG.Wait()
}

func (r *replica) GetMasterConn() net.Conn {
	return r.masterConn
}

func (r *replica) acceptClientConnections() {
	address := fmt.Sprintf("%s:%d", r.args.Host, r.args.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to bind to address: %s\n", address)
	}
	defer listener.Close()

	r.acceptClientsReady <- struct{}{}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go r.handleClient(nil, conn, true)
	}
}

func (r *replica) connectToMaster() {
	r.dialMaster()
	r.processMasterHandshake()
	r.handleMaster()
}

func (r *replica) dialMaster() {
	address := fmt.Sprintf("%s:%d", r.args.ReplicaOf.Host, r.args.ReplicaOf.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to dial master address: %s\n: %v", address, err)
	}
	r.masterConn = conn
}

func (r *replica) processMasterHandshake() {
	r.processMasterHandshakePING()
	r.processMasterHandshakeREPLCONF()
	r.processMasterHandshakePSYNC()
}

func (r *replica) processMasterHandshakePING() {
	pingCommand := resp.CreateBulkStringArray("PING")
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
		resp.CreateBulkStringArray("REPLCONF", "listening-port", strconv.Itoa(r.args.Port)),
		resp.CreateBulkStringArray("REPLCONF", "capa", "psync2"),
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
	psyncCommand := resp.CreateBulkStringArray("PSYNC", "?", "-1")
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
		log.Fatalf("Master handshake PSYNC (3/3) psync response has wrong RESP type, expected: %s, got %T", reflect.TypeOf(resp.SimpleString{}).String(), simpleString)
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

	rdbPayload, restBytes, err := r.readRDBFileFromMaster()
	if err != nil {
		log.Fatalf("Master handshake PSYNC (3/3) read RDB file from master error: %s\n", err)
	}
	r.processRDBFile(rdbPayload)

	if len(restBytes) > 0 {
		r.masterConnBuffer = r.processCommands(restBytes, r.masterConn, false)
	} else {
		r.masterConnBuffer = nil
	}
}

func (r *replica) handleMaster() {
	r.handleClient(r.masterConnBuffer, r.masterConn, false)
}

func (r *replica) readFromMaster() ([]byte, int, error) {
	b := make([]byte, 1024)
	n, err := r.masterConn.Read(b)
	if errors.Is(err, io.EOF) {
		return b, n, nil
	}

	return b, n, err
}

func (r *replica) readValueFromMaster() (resp.Value, error) {
	if len(r.masterConnBuffer) == 0 {
		b, n, err := r.readFromMaster()
		if err != nil {
			return nil, fmt.Errorf("read from master error: %v", err)
		}
		r.masterConnBuffer = append(r.masterConnBuffer, b[:n]...)
	}

	rest, value, err := r.respController.Decode(r.masterConnBuffer)
	if err != nil {
		return nil, fmt.Errorf("RESP controller decode error: %v", err)
	}

	r.masterConnBuffer = rest

	return value, nil
}

func (r *replica) readRDBFileFromMaster() ([]byte, []byte, error) {
	b := r.masterConnBuffer
	n := len(b)

	b, n, err := r.ensureInitialRDBData(b, n)
	if err != nil {
		return nil, nil, err
	}

	if n == 0 || b[0] != '$' {
		r.masterConnBuffer = b
		return nil, b, nil
	}

	b, n, i, err := r.findLengthDelimiter(b, n)
	if err != nil {
		return nil, nil, err
	}

	rdbFileLen, err := r.parseRDBFileLength(b, i)
	if err != nil {
		return nil, nil, err
	}

	fileContentsIdx := i + 2
	b, _, err = r.readFullRDBFileContent(b, n, fileContentsIdx, rdbFileLen)
	if err != nil {
		return nil, nil, err
	}

	rdbPayload := b[fileContentsIdx : fileContentsIdx+rdbFileLen]
	restBytes := b[fileContentsIdx+rdbFileLen:]

	r.masterConnBuffer = restBytes

	return rdbPayload, restBytes, nil
}

func (r *replica) ensureInitialRDBData(b []byte, n int) ([]byte, int, error) {
	if n == 0 || n < 2 || b[0] != '$' {
		tmp, nRead, err := r.readFromMaster()
		if err != nil {
			return nil, 0, err
		}
		b = append(b, tmp[:nRead]...)
		n = len(b)
	}
	return b, n, nil
}

func (r *replica) findLengthDelimiter(b []byte, n int) ([]byte, int, int, error) {
	i := bytes.Index(b, []byte("\r\n"))
	if i == -1 {
		for i == -1 && n < 4096 {
			var err error
			b, err = r.appendDataFromMaster(b)
			if err != nil {
				return nil, 0, 0, err
			}
			n = len(b)
			i = bytes.Index(b, []byte("\r\n"))
		}
		if i == -1 {
			return nil, 0, 0, fmt.Errorf("could not find end of RDB file length")
		}
	}
	return b, n, i, nil
}

func (r *replica) parseRDBFileLength(b []byte, i int) (int, error) {
	rdbFileLen, err := strconv.Atoi(string(b[1:i]))
	if err != nil {
		return 0, fmt.Errorf("RDB file length parse error: %v", err)
	}
	return rdbFileLen, nil
}

func (r *replica) readFullRDBFileContent(b []byte, n, fileContentsIdx, rdbFileLen int) ([]byte, int, error) {
	if fileContentsIdx+rdbFileLen > n {
		for fileContentsIdx+rdbFileLen > n {
			var err error
			b, err = r.appendDataFromMaster(b)
			if err != nil {
				return nil, 0, err
			}
			n = len(b)
		}
	}
	return b, n, nil
}

func (r *replica) appendDataFromMaster(b []byte) ([]byte, error) {
	tmp, nRead, err := r.readFromMaster()
	if err != nil {
		return nil, err
	}
	if nRead == 0 {
		return nil, fmt.Errorf("unexpected EOF while reading data from master")
	}
	b = append(b, tmp[:nRead]...)
	return b, nil
}

func (*replica) initReplicationInfo() *replication.Info {
	return &replication.Info{
		Role:             "slave",
		MasterReplID:     "?",
		MasterReplOffset: -1,
	}
}
