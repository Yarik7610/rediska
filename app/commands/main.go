package commands

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Controller struct {
	storage     *memory.Storage
	args        *config.Args
	replication replication.Base
}

func NewController(storage *memory.Storage, args *config.Args, replication replication.Base) *Controller {
	return &Controller{storage: storage, args: args, replication: replication}
}

func (c *Controller) HandleCommand(cmd resp.Value, conn net.Conn, writeResponseToConn bool) error {
	result := c.handleCommand(cmd, conn)
	if writeResponseToConn && result != nil {
		err := c.Write(result, conn)
		if err != nil {
			return err
		}
	}
	c.updateMasterReplOffset(cmd, conn)
	return nil
}

func (c *Controller) Write(cmd resp.Value, conn net.Conn) error {
	encoded, err := cmd.Encode()
	if err != nil {
		fmt.Fprintf(conn, "-ERR encode error: %v\r\n", err)
		return err
	}
	_, err = conn.Write(encoded)
	return err
}

func (c *Controller) updateMasterReplOffset(cmd resp.Value, conn net.Conn) error {
	b, err := cmd.Encode()
	if err != nil {
		return err
	}
	l := len(b)

	if r, ok := c.replication.(replication.Replica); ok {
		if r.GetMasterConn() == conn {
			c.replication.IncrMasterReplOffset(l)
		}
	}
	return nil
}

func (c *Controller) handleCommand(cmd resp.Value, conn net.Conn) resp.Value {
	switch cmd := cmd.(type) {
	case resp.Array:
		return c.handleArrayCommand(cmd, conn)
	case resp.SimpleString:
		return c.handleSimpleStringCommand(cmd)
	default:
		return resp.SimpleError{Value: "commands must be sent as RESP array or simple string"}
	}
}

func (c *Controller) handleArrayCommand(cmd resp.Array, conn net.Conn) resp.Value {
	if len(cmd.Value) == 0 {
		return resp.SimpleError{Value: "empty RESP command array"}
	}

	commandAndArgs, err := extractCommandAndArgs(cmd.Value)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("extract command and args from RESP array command error: %v", err)}
	}

	command := commandAndArgs[0]
	args := commandAndArgs[1:]
	switch strings.ToUpper(command) {
	case "PING":
		return c.ping()
	case "ECHO":
		return c.echo(args)
	case "GET":
		return c.get(args)
	case "SET":
		return c.set(args, commandAndArgs)
	case "CONFIG":
		secondCommand := strings.ToUpper(args[0])
		if secondCommand == "GET" {
			return c.configGet(args[1:])
		}
		return resp.SimpleError{Value: fmt.Sprintf("unknown command CONFIG '%s'", secondCommand)}
	case "KEYS":
		return c.keys(args)
	case "INFO":
		return c.info(args)
	case "REPLCONF":
		return c.replconf(args, conn)
	case "PSYNC":
		return c.psync(args, conn)
	case "WAIT":
		return c.wait(args)
	case "DEL":
		return c.del(args, commandAndArgs)
	case "RPUSH":
		return c.rpush(args, commandAndArgs)
	case "LPUSH":
		return c.lpush(args, commandAndArgs)
	case "LPOP":
		return c.lpop(args, commandAndArgs)
	case "RPOP":
		return c.rpop(args, commandAndArgs)
	case "LRANGE":
		return c.lrange(args)
	case "LLEN":
		return c.llen(args)
	default:
		return resp.SimpleError{Value: fmt.Sprintf("unknown command '%s'", command)}
	}
}

func (c *Controller) handleSimpleStringCommand(cmd resp.SimpleString) resp.Value {
	switch cmd.Value {
	case "PING":
		return c.ping()
	default:
		return resp.SimpleError{Value: fmt.Sprintf("unknown command '%s'", cmd.Value)}
	}
}

func (c *Controller) propagateWriteCommand(commandAndArgs []string) {
	if m, ok := c.replication.(replication.Master); ok {
		m.SetHasPendingWrites(true)
		go m.Propagate(commandAndArgs)
	}
}
