package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/state"
)

func main() {
	args := config.NewArgs()

	server := state.SpawnServer(args)
	server.Start()
}
