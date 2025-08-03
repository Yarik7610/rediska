package utils

import "net"

func GetRemoteAddr(conn net.Conn) string {
	return conn.RemoteAddr().String()
}
