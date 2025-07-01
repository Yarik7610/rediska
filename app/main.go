package main

import (
	"log"
	"net"
)

func handleClient(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("+PONG\r\n"))
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatalln("Failed to bind to '0.0.0.0:6379'")
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go handleClient(conn)
	}
}
