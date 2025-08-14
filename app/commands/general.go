package commands

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) ping(conn net.Conn) resp.Value {
	subscribeModePong := "pong"
	subscribeModeEmptyStr := ""

	if c.pubsubController.InSubscribeMode(conn) {
		return resp.Array{Value: []resp.Value{resp.BulkString{Value: &subscribeModePong}, resp.BulkString{Value: &subscribeModeEmptyStr}}}
	}
	return resp.SimpleString{Value: "PONG"}
}

func (*controller) echo(args []string) resp.BulkString {
	res := strings.Join(args, " ")
	return resp.BulkString{Value: &res}
}

func (c *controller) keys(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "KEYS command error: only 1 argument supported"}
	}

	pattern := args[0]
	if pattern != "*" {
		return resp.SimpleError{Value: "KEYS command error: only '*' pattern supported"}
	}

	keys := c.storage.Keys()
	return resp.CreateBulkStringArray(keys...)
}

func (c *controller) del(args, commandAndArgs []string) resp.Value {
	for _, key := range args {
		c.storage.Del(key)
	}

	c.propagateWriteCommand(commandAndArgs)
	return resp.SimpleString{Value: "OK"}
}

func (c *controller) configGet(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "CONFIG GET command must have only 1 arg"}
	}

	arg := args[0]
	value := []string{arg}
	switch arg {
	case "host":
		value = append(value, c.args.Host)
	case "port":
		value = append(value, strconv.Itoa(c.args.Port))
	case "replicaof":
		var replicaOfString string
		if c.args.ReplicaOf != nil {
			replicaOfString = c.args.ReplicaOf.String()
		}
		value = append(value, replicaOfString)
	case "dir":
		value = append(value, c.args.DBDir)
	case "dbfilename":
		value = append(value, c.args.DBFilename)
	default:
		return resp.SimpleError{Value: fmt.Sprintf("CONFIG GET command unknown arg: %s", arg)}
	}

	return resp.CreateBulkStringArray(value...)
}

func (c *controller) valuetype(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "TYPE command must have only 1 arg"}
	}

	key := args[0]

	return resp.SimpleString{Value: c.storage.Type(key)}
}
