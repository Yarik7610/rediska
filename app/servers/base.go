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
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type base struct {
	Storage           *memory.Storage
	RESPController    *resp.Controller
	CommandController *commands.Controller
	Args              *config.Args
}

func newBase(args *config.Args) *base {
	return &base{
		Storage:        memory.NewStorage(),
		RESPController: resp.NewController(),
		Args:           args,
	}
}

func (base *base) initStorage() {
	if base.Args.DBDir == "" || base.Args.DBFilename == "" {
		return
	}

	if !rdb.IsFileExists(base.Args.DBDir, base.Args.DBFilename) {
		return
	}

	b, err := rdb.ReadRDB(base.Args.DBDir, base.Args.DBFilename)
	if err != nil {
		log.Printf("Skip RDB storage seed, RDB file read error: %v\n", err)
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
				base.Storage.SetWithExpiry(key, item.Value, item.Expires)
			}
		} else {
			base.Storage.Set(key, item.Value)
		}
	}
}

func (base *base) acceptClientConnections() {
	address := fmt.Sprintf("%s:%d", base.Args.Host, base.Args.Port)

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

		for len(buf) > 0 {
			rest, value, err := base.RESPController.Decode(buf)
			if err != nil {
				log.Printf("decode error: %v", err)
				fmt.Fprintf(conn, "-ERR %v\r\n", err)
				return
			}

			buf = rest

			base.CommandController.HandleCommand(value, conn)
		}
	}
}

func (base *base) startExpiredKeysCleanup() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			base.Storage.CleanExpiredKeys()
		}
	}()
}
