package commands

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/pubsub"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/transaction"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

type Controller interface {
	HandleCommand(cmd resp.Value, conn net.Conn, writeResponseToConn bool) (resp.Value, error)
}

type controller struct {
	args                  *config.Args
	storage               memory.MultiTypeStorage
	replicationController replication.BaseController
	pubsubController      pubsub.Controller
	transactionController transaction.Controller
}

func NewController(
	args *config.Args,
	storage memory.MultiTypeStorage,
	replicationController replication.BaseController,
	pubsubController pubsub.Controller,
	transactionController transaction.Controller) Controller {
	return &controller{
		args:                  args,
		storage:               storage,
		replicationController: replicationController,
		pubsubController:      pubsubController,
		transactionController: transactionController,
	}
}

func (c *controller) HandleCommand(cmd resp.Value, conn net.Conn, writeResponseToConn bool) (resp.Value, error) {
	result := c.handleCommand(cmd, conn)
	if writeResponseToConn && result != nil {
		err := utils.WriteCommand(result, conn)
		if err != nil {
			return nil, err
		}
	}
	c.updateMasterReplOffset(cmd, conn)
	return result, nil
}

func (c *controller) updateMasterReplOffset(cmd resp.Value, conn net.Conn) error {
	if r, ok := c.replicationController.(replication.ReplicaController); ok && r.GetMasterConn() == conn {
		b, err := cmd.Encode()
		if err != nil {
			return err
		}
		l := len(b)
		c.replicationController.IncrMasterReplOffset(l)
	}
	return nil
}

func (c *controller) handleCommand(cmd resp.Value, conn net.Conn) resp.Value {
	switch cmd := cmd.(type) {
	case resp.Array:
		return c.handleArrayCommand(cmd, conn)
	case resp.SimpleString:
		return c.handleSimpleStringCommand(cmd, conn)
	default:
		return resp.SimpleError{Value: "commands must be sent as RESP array or simple string"}
	}
}

func (c *controller) handleArrayCommand(cmd resp.Array, conn net.Conn) resp.Value {
	if len(cmd.Value) == 0 {
		return resp.SimpleError{Value: "empty RESP command array"}
	}

	commandAndArgs, err := extractCommandAndArgs(cmd.Value)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("extract command and args from RESP array command error: %v", err)}
	}

	command := commandAndArgs[0]
	args := commandAndArgs[1:]

	if c.transactionController.InTransaction(conn) && !c.transactionController.IsTransactionCommand(command) {
		// In original redis logic is a bit more complicated
		// First all commands are validated and then they are executed or queued
		// In mine, all commands, even invalid, are queued and when EXEC works i response with resp.SimpleError if there is an error
		c.transactionController.EnqueueCommand(conn, cmd)
		return resp.SimpleString{Value: "QUEUED"}
	}

	if c.pubsubController.InSubscribeMode(conn) && !c.pubsubController.IsSubscribeModeCommand(command) {
		return resp.SimpleError{
			Value: fmt.Sprintf("ERR Can't execute '%s': only (P|S)SUBSCRIBE / (P|S)UNSUBSCRIBE / PING / QUIT / RESET are allowed in this context", strings.ToLower(command)),
		}
	}

	switch strings.ToUpper(command) {
	case "PING":
		return c.ping(conn)
	case "ECHO":
		return c.echo(args)
	case "GET":
		return c.get(args)
	case "INCR":
		return c.incr(args)
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
	case "RPUSH", "LPUSH":
		return c.push(commandAndArgs)
	case "RPOP", "LPOP":
		return c.pop(commandAndArgs)
	case "BRPOP", "BLPOP":
		return c.bpop(commandAndArgs)
	case "LRANGE":
		return c.lrange(args)
	case "LLEN":
		return c.llen(args)
	case "TYPE":
		return c.valuetype(args)
	case "XADD":
		return c.xadd(args, commandAndArgs)
	case "XRANGE":
		return c.xrange(args)
	case "XREAD":
		return c.xread(args)
	case "SUBSCRIBE", "UNSUBSCRIBE":
		return c.subcribeOrUnsubscribe(commandAndArgs, conn)
	case "PUBLISH":
		return c.publish(args)
	case "MULTI":
		return c.multi(args, conn)
	case "EXEC":
		return c.exec(args, conn)
	case "DISCARD":
		return c.discard(args, conn)
	case "ZADD":
		return c.zadd(args, commandAndArgs)
	case "ZRANK":
		return c.zrank(args)
	default:
		return resp.SimpleError{Value: fmt.Sprintf("unknown command '%s'", command)}
	}
}

func (c *controller) handleSimpleStringCommand(cmd resp.SimpleString, conn net.Conn) resp.Value {
	switch cmd.Value {
	case "PING":
		return c.ping(conn)
	default:
		return resp.SimpleError{Value: fmt.Sprintf("unknown command '%s'", cmd.Value)}
	}
}

func (c *controller) propagateWriteCommand(commandAndArgs []string) {
	if m, ok := c.replicationController.(replication.MasterController); ok {
		m.SetHasPendingWrites(true)
		go m.Propagate(commandAndArgs)
	}
}
