package commands

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func HandleCommand(conn net.Conn, unit resp.Value) {
	response := handleCommand(unit)
	encoded, err := response.Encode()
	if err != nil {
		fmt.Fprintf(conn, "-ERR encode error: %v\r\n", err)
		return
	}
	conn.Write(encoded)
}

func handleCommand(unit resp.Value) resp.Value {
	switch u := unit.(type) {
	case resp.Array:
		return handleArrayCommand(u)
	case resp.SimpleString:
		return handleSimpleStringCommand(u)
	default:
		return resp.SimpleError{Value: "commands must be sent as RESP array or simple string"}
	}
}

func handleArrayCommand(unit resp.Array) resp.Value {
	if len(unit.Value) == 0 {
		return resp.SimpleError{Value: "empty RESP array"}
	}

	commandAndArgs, err := extractCommandAndArgs(unit.Value)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("extract command and args from RESP array error: %v", err)}
	}

	command := commandAndArgs[0]

	switch strings.ToUpper(command) {
	case "PING":
		return Ping()
	case "ECHO":
		return Echo(commandAndArgs[1:])
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
