package main

import (
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/state"
)

func main() {
	args := config.NewArgs()

	address := fmt.Sprintf("%s:%d", args.Host, args.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to bind to adress %s\n", address)
	}
	defer listener.Close()

	server := state.SpawnServer(args, listener)
	server.Start()
}
