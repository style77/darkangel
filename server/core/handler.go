package core

import (
	"bytes"
	"fmt"
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

	if len(request.Args) != len(cmd.Arguments) {
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

func execBashCmd(request Request) CommandResult {
	target := request.Server.GetConnection(request.Args[0])
	if target == nil {
		return HandleWrongTarget(request)
	}
	fmt.Println(target)

	script := request.Args[2] // todo execute this script on target machine and return response with channels

	return CommandResult{}
}
