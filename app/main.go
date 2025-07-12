package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/state"
)

func handleClient(conn net.Conn, server state.Server) {
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

		value, err := server.DecodeRESP(b[:n])
		if err != nil {
			fmt.Fprintf(conn, "RESP controller decode error: %v\n", err)
			continue
		}

		server.HandleCommand(value, conn)
	}
}

func main() {
	args := config.NewArgs()

	address := fmt.Sprintf("%s:%d", args.Host, args.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to bind to adress %s\n", address)
	}
	defer listener.Close()

	server := state.SpawnServer(args)
	server.Start()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go handleClient(conn, server)
	}
}
