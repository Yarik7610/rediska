package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

func handleClient(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		text := scanner.Text()
		if strings.TrimSpace(text) == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		}
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go handleClient(conn)
	}
}
