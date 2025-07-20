package servers

import (
	"errors"
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
	if b == nil {
		log.Println("Skip RDB storage seed, no RDB file detected")
		return
	}

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

func (base *base) handleClient(initialBuffer []byte, conn net.Conn, writeResponseToConn bool) {
	defer conn.Close()

	buf := make([]byte, 0, 4096)
	if initialBuffer != nil {
		log.Println("INITIAL BUFFER", string(initialBuffer))
		buf = append(buf, initialBuffer...)
	}
	tmp := make([]byte, 1024)

	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("connection %s closed: (EOF)", conn.RemoteAddr().String())
				return
			}
			log.Printf("read error: %v", err)
			return
		}
		buf = append(buf, tmp[:n]...)
		buf = base.processCommands(buf, conn, writeResponseToConn)
	}
}

func (base *base) processCommands(buf []byte, conn net.Conn, writeResponseToConn bool) []byte {
	for len(buf) > 0 {
		rest, value, err := base.respController.Decode(buf)
		if err != nil {
			log.Printf("decode error: %v", err)
			return buf
		}

		err = base.commandController.HandleCommand(value, conn, writeResponseToConn)
		if err != nil {
			log.Printf("handle command error: %v, continue to work", err)
		}

		buf = rest
	}
	return buf
}

func (base *base) startExpiredKeysCleanup() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			base.storage.CleanExpiredKeys()
		}
	}()
}
