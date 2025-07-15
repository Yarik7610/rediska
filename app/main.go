package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/servers"
)

func main() {
	args := config.NewArgs()

	server := servers.SpawnServer(args)
	server.Start()
}
