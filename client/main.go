package main

import (
	"fmt"
	"log"
	"net"
)

const host = "localhost"
const port = 2306

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal("Couldn't connect to server")
	}

	log.Printf("Connected to server")

	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		go HandleCommand(buffer)
	}
}
