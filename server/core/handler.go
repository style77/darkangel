package core

import (
	"bytes"
	"darkangel/server/constant"
	"encoding/json"
	"fmt"
	"log"
	"text/tabwriter"
)

func HandleCommand(request Request) CommandResult {
	var cmd Command
	exists := false

	if request.CommandName == "help" {
		cmd, exists = helpCmd, true
	} else {
		cmd, exists = commandsMap[request.CommandName]
	}
	if !exists {
		return HandleCommandDoesNotExists(request.CommandName)
	}

	if len(request.Args) < len(cmd.Arguments) {
		return HandleCommandWrongArgumentsLength(cmd, request)
	}

	return cmd.Callback(request)
}

func connectionsCmd(request Request) CommandResult {
	conns := request.Server.GetConnections()

	var buffer bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buffer, 0, 8, 0, '\t', 0)

	fmt.Fprintf(w, "Local Address\tRemote Address\t(%d connections)\n", len(conns))
	fmt.Fprintln(w, "--------------\t--------------\t")

	for _, client := range conns {
		fmt.Fprintf(w, "%s\t%s\t\n", client.Conn.LocalAddr(), client.Conn.RemoteAddr())
	}

	w.Flush()

	formattedTable := buffer.String()

	return CommandResult{
		CommandName: request.CommandName,
		Result:      true,
		Output:      formattedTable,
	}
}

func packData(data map[string]interface{}) []byte {
	v, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Couldn't Marshall data")
	}
	return v
}

func execPsCmd(request Request) CommandResult {
	targets := request.Server.GetConnection(request.Args[0])
	if targets == nil {
		return HandleWrongTarget(request)
	}

	for _, target := range targets {
		data := packData(map[string]interface{}{"name": "ps", "args": request.Args[1:]})
		_, err := target.Conn.Write([]byte(data))
		if err != nil {
			log.Fatalln(fmt.Sprintf(constant.Red+"Couldn't write to target %s"+constant.Reset, target.Conn.RemoteAddr()))
		}
	}

	return CommandResult{}
}
