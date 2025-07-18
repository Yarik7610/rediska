package servers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/persistence/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type base struct {
	storage           *memory.Storage
	respController    *resp.Controller
	commandController *commands.Controller
	args              *config.Args
	replicationInfo   *replication.Info
}

var _ replication.Base = (*base)(nil)

func newBase(args *config.Args) *base {
	return &base{
		storage:        memory.NewStorage(),
		respController: resp.NewController(),
		args:           args,
	}
}

func (base *base) Info() *replication.Info {
	return base.replicationInfo
}

func (base *base) SetMasterReplID(replID string) {
	base.replicationInfo.MasterReplID = replID
}

func (base *base) SetMasterReplOfffset(replOffset int) {
	base.replicationInfo.MasterReplOffset = replOffset
}

func (base *base) initStorage() {
	if base.args.DBDir == "" || base.args.DBFilename == "" {
		return
	}

	if !rdb.IsFileExists(base.args.DBDir, base.args.DBFilename) {
		return
	}
	base.persistWithRDBFile()
}

func (base *base) persistWithRDBFile() {
	b, err := rdb.ReadRDBFile(base.args.DBDir, base.args.DBFilename)
	if err != nil {
		log.Printf("Skip RDB storage seed, RDB file read error: %v\n", err)
		return
	}
	base.processRDBFile(b)
}

func (base *base) processRDBFile(b []byte) {
	items, err := rdb.Decode(b)
	if err != nil {
		log.Printf("Skip RDB storage seed, RDB decode error: %v\n", err)
		return
	}
	base.putRDBItemsIntoStorage(items)
}

func (base *base) putRDBItemsIntoStorage(items map[string]memory.Item) {
	for key, item := range items {
		if memory.ItemHasExpiration(&item) {
			if !memory.ItemExpired(&item) {
				base.storage.SetWithExpiry(key, item.Value, item.Expires)
			}
		} else {
			base.storage.Set(key, item.Value)
		}
	}
}

func (base *base) acceptClientConnections() {
	address := fmt.Sprintf("%s:%d", base.args.Host, base.args.Port)

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

		go base.handleClient(conn)
	}
}

func (base *base) handleClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 1024)

	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Printf("read error: %v", err)
			return
		}
		buf = append(buf, tmp[:n]...)

		//Reading all commands from buffer, if there is more than 1 command
		for len(buf) > 0 {
			rest, value, err := base.respController.Decode(buf)
			if err != nil {
				log.Printf("decode error: %v", err)
				fmt.Fprintf(conn, "-ERR %v\r\n", err)
				return
			}

			buf = rest

			base.commandController.HandleCommand(value, conn)
		}
	}
}

func (base *base) startExpiredKeysCleanup() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			base.storage.CleanExpiredKeys()
		}
	}()
}
