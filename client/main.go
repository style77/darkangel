package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

const host = "localhost"
const port = 2306

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal("Couldn't connect to server")
	}

	defer conn.Close()
	log.Printf("Connected to server")

	for {
		reply := make([]byte, 2048)

		_, err = conn.Read(reply)
		if err != nil {
			log.Println("Write to server failed:", err.Error())
			os.Exit(1)
		}

		reply = bytes.Trim(reply, "\x00")
		log.Println("Server send data:", string(reply))

		// if server send A, this means we are only checking connection
		if string(reply) == "A" {
			continue
		}

		data := Command{}
		err := json.Unmarshal(reply, &data)
		if err != nil {
			log.Printf("Couldn't unpack data. %s", err)
			continue
		}

		HandleCommand(Command{
			Name: data.Name,
			Args: data.Args,
		}, conn)
	}
}

type Output struct {
	Output     string
	RemoteAddr string
}

func sendOutputToServer(output string, conn net.Conn) {
	var buffer bytes.Buffer

	enc := gob.NewEncoder(&buffer)

	err := enc.Encode(Output{
		Output:     output,
		RemoteAddr: conn.RemoteAddr().String(),
	})
	if err != nil {
		log.Println("Couldn't encode output")
	}
	conn.Write(buffer.Bytes())
}
