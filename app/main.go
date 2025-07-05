package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/state"
)

func handleClient(conn net.Conn, server *state.Server) {
	defer conn.Close()

	for {
		b := make([]byte, 1024)

		n, err := conn.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Printf("Connnection read error: %v", err)
			return
		}

		value, err := server.RESPController.Decode(b[:n])
		if err != nil {
			fmt.Fprintf(conn, "RESP controller decode error: %v\n", err)
			continue
		}

		commands.HandleCommand(value, conn, server)
	}
}

func main() {
	port := flag.Int("port", 6379, "The port of redis server")
	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to bind to adress %s\n", address)
	}
	defer listener.Close()

	server := state.NewServer()
	server.StartExpiredKeysCleanup()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go handleClient(conn, server)
	}
}
