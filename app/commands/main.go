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
	storage         *memory.Storage
	args            *config.Args
	replicationInfo *replication.Info
}

func NewController(storage *memory.Storage, args *config.Args, replicationInfo *replication.Info) *Controller {
	return &Controller{storage: storage, args: args, replicationInfo: replicationInfo}
}

func (c *Controller) HandleCommand(unit resp.Value, conn net.Conn) {
	response := c.handleCommand(unit)
	encoded, err := response.Encode()
	if err != nil {
		fmt.Fprintf(conn, "-ERR encode error: %v\r\n", err)
		return
	}
	conn.Write(encoded)
}

func (c *Controller) handleCommand(unit resp.Value) resp.Value {
	switch u := unit.(type) {
	case resp.Array:
		return c.handleArrayCommand(u)
	case resp.SimpleString:
		return c.handleSimpleStringCommand(u)
	default:
		return resp.SimpleError{Value: "commands must be sent as RESP array or simple string"}
	}
}

func (c *Controller) handleArrayCommand(unit resp.Array) resp.Value {
	if len(unit.Value) == 0 {
		return resp.SimpleError{Value: "empty RESP array"}
	}

	commandAndArgs, err := extractCommandAndArgs(unit.Value)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("extract command and args from RESP array error: %v", err)}
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
		return c.set(args)
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
	default:
		return resp.SimpleError{Value: fmt.Sprintf("unknown command '%s'", command)}
	}
}

func (c *Controller) handleSimpleStringCommand(unit resp.SimpleString) resp.Value {
	switch unit.Value {
	case "PING":
		return c.ping()
	default:
		return resp.SimpleError{Value: fmt.Sprintf("unknown command '%s'", unit.Value)}
	}
}
