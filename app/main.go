package main

import (
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatalln("Failed to bind to '0.0.0.0:6379'")
	}

	_, err = listener.Accept()
	if err != nil {
		log.Fatalf("Error accepting connection: %v\n", err)
	}
}
