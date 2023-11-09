package client

import (
	"encoding/gob"
	"fmt"
	"net"
)

type Client struct {
	Conn net.Conn
}

type Output struct {
	Output     string
	RemoteAddr string
}

func (client *Client) HandleRequest() {
	dec := gob.NewDecoder(client.Conn)
	for {
		var output Output
		if err := dec.Decode(&output); err != nil {
			client.Conn.Close()
			return
		}
		fmt.Printf("\n\rOutput from %s: %s\n", output.RemoteAddr, output.Output)
	}
}
