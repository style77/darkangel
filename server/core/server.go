package core

import (
	"bufio"
	"darkangel/server/client"
	"darkangel/server/constant"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Server struct {
	host        string
	port        string
	connections []*client.Client
}

func (server *Server) ListenForInput() {
	newlineCounter := 0
	for {
		fmt.Print(strings.Repeat("\n", newlineCounter) + "> ")
		newlineCounter = 0 // reset counter
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		formattedInput := strings.ToLower(input)
		splittedInput := strings.Fields(formattedInput)

		if len(splittedInput) == 0 {
			continue
		}

		commandName := splittedInput[0]
		args := splittedInput[1:]

		result := HandleCommand(Request{
			Server:      server,
			CommandName: commandName,
			Args:        args,
		})

		fmt.Print(result.Output)
		newlineCounter = 2
	}
}

func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(constant.Yellow+"Server listening on %s:%s"+constant.Reset, string(server.host), string(server.port))

	go server.ListenForInput()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(constant.Yellow+"User %s connected"+constant.Reset, conn.RemoteAddr())

		client := &client.Client{
			Conn: conn,
		}
		server.connections = append(server.connections, client)
		go client.HandleRequest()
	}
}

func (server *Server) GetConnections() []*client.Client {
	for n, connectedClient := range server.connections {
		// Send 1 byte to ensure that client is still active
		_, err := connectedClient.Conn.Write([]byte{65})
		if err != nil {
			server.connections = append(server.connections[:n], server.connections[n+1:]...)
		}
	}
	return server.connections
}

func (server *Server) GetConnection(target string) []*client.Client {
	if target == "all" {
		return server.GetConnections()
	}

	for _, connectedClient := range server.connections {
		if connectedClient.Conn.RemoteAddr().String() == target {
			_, err := connectedClient.Conn.Write([]byte{65})
			if err != nil {
				return nil
			}

			return []*client.Client{connectedClient}
		}
	}

	return nil
}
