package client

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	Conn net.Conn
}

func (client *Client) HandleRequest() {
	reader := bufio.NewReader(client.Conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			client.Conn.Close()
			return
		}
		fmt.Printf("Message incoming: %s", string(message))
		client.Conn.Write([]byte("Message received.\n"))
	}
}
