package utils

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func WriteCommand(cmd resp.Value, conn net.Conn) error {
	encoded, err := cmd.Encode()
	if err != nil {
		return err
	}
	_, err = conn.Write(encoded)
	return err
}
