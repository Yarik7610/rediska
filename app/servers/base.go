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
	"github.com/codecrafters-io/redis-starter-go/app/pubsub"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/transaction"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

type base struct {
	args                  *config.Args
	storage               memory.MultiTypeStorage
	respController        resp.Controller
	pubsubController      pubsub.Controller
	transactionController transaction.Controller
	commandController     commands.Controller
}

func newBase(args *config.Args) *base {
	return &base{
		args:                  args,
		storage:               memory.NewMultiTypeStorage(),
		respController:        resp.NewController(),
		pubsubController:      pubsub.NewController(),
		transactionController: transaction.NewController(),
	}
}

func (base *base) initStorage() {
	if base.args.DBDir == "" || base.args.DBFilename == "" {
		return
	}

	if !rdb.IsFileExists(base.args.DBDir, base.args.DBFilename) {
		log.Printf("Skip RDB storage seed, no such file or directory found")
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
	base.putRDBStringItemsIntoStorage(items)
}

func (base *base) putRDBStringItemsIntoStorage(items map[string]memory.String) {
	for key, item := range items {
		if base.storage.StringStorage().ItemHasExpiration(&item) {
			if !base.storage.StringStorage().ItemExpired(&item) {
				base.storage.StringStorage().SetWithExpiry(key, item.Value, item.Expires)
			}
		} else {
			base.storage.StringStorage().Set(key, item.Value)
		}
	}
}

func (base *base) listenTCP() net.Listener {
	address := fmt.Sprintf("%s:%d", base.args.Host, base.args.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to bind to address: %s\n", address)
	}
	return listener
}

func (base *base) handleClient(initialBuffer []byte, conn net.Conn, writeResponseToConn bool) {
	defer conn.Close()

	buf := make([]byte, 0, 4096)
	if initialBuffer != nil {
		buf = append(buf, initialBuffer...)
	}
	tmp := make([]byte, 1024)

	for {
		n, err := conn.Read(tmp)
		if err != nil {
			addr := utils.GetRemoteAddr(conn)
			if errors.Is(err, io.EOF) {
				log.Printf("Connection %s closed: (EOF)", addr)
				base.pubsubController.UnsubscribeFromAllChannels(conn)
				return
			}
			log.Printf("Connection %s read error: %v", addr, err)
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
			log.Printf("RESP controller decode error: %v", err)
			return buf
		}

		_, err = base.commandController.HandleCommand(value, conn, writeResponseToConn)
		if err != nil {
			log.Printf("handle command error: %v, continue to work", err)
		}

		buf = rest
	}
	return buf
}

func (base *base) startExpiredStringKeysCleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		base.storage.StringStorage().CleanExpiredKeys()
	}
}
