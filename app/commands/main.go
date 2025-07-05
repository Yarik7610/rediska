package commands

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/state"
)

func HandleCommand(unit resp.Value, conn net.Conn, server *state.Server) {
	response := handleCommand(unit, server)
	encoded, err := response.Encode()
	if err != nil {
		fmt.Fprintf(conn, "-ERR encode error: %v\r\n", err)
		return
	}
	conn.Write(encoded)
}

func handleCommand(unit resp.Value, server *state.Server) resp.Value {
	switch u := unit.(type) {
	case resp.Array:
		return handleArrayCommand(u, server)
	case resp.SimpleString:
		return handleSimpleStringCommand(u)
	default:
		return resp.SimpleError{Value: "commands must be sent as RESP array or simple string"}
	}
}

func handleArrayCommand(unit resp.Array, server *state.Server) resp.Value {
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
		return Ping()
	case "ECHO":
		return Echo(args)
	case "GET":
		return Get(args, server)
	case "SET":
		return Set(args, server)
	case "CONFIG":
		secondCommand := strings.ToUpper(args[0])
		if secondCommand == "GET" {
			return ConfigGet(args[1:], server)
		}
		return resp.SimpleError{Value: fmt.Sprintf("unknown command CONFIG '%s'", secondCommand)}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("unknown command '%s'", command)}
	}
}

func handleSimpleStringCommand(unit resp.SimpleString) resp.Value {
	switch unit.Value {
	case "PING":
		return Ping()
	default:
		return resp.SimpleError{Value: fmt.Sprintf("unknown command '%s'", unit.Value)}
	}
}
