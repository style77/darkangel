package core

import (
	"fmt"
)

type CommandResult struct {
	CommandName string
	Result      bool // True completed, False error/uncompleted
	Output      string
}

type Argument struct {
	Name        string
	Description string
}

type CallbackFunc func(request Request) CommandResult
type Command struct {
	Name        string
	Description string
	Callback    CallbackFunc
	Arguments   []Argument
}

var commandsMap = map[string]Command{
	"connections": {
		Name:        "connections",
		Description: "Display connections",
		Callback:    connectionsCmd,
		Arguments:   []Argument{},
	},
	"execbash": {
		Name:        "execbash",
		Description: "Run bash script",
		Callback:    execBashCmd,
		Arguments: []Argument{
			{Name: "target", Description: "target remote address"},
			{Name: "script", Description: "bash script to execute"},
		},
	},
}

var helpCmd = Command{
	Name:        "help",
	Description: "Display available commands",
	Callback: func(request Request) CommandResult {
		helpMessage := generateHelpMessage()
		return CommandResult{Output: helpMessage}
	},
}

func generateHelpMessage() string {
	helpMessage := "Available commands:\n"

	for _, cmd := range commandsMap {
		helpMessage += fmt.Sprintf("%s - %s\n", cmd.Name, cmd.Description)

		if len(cmd.Arguments) > 0 {
			helpMessage += "  Arguments:\n"
			for _, arg := range cmd.Arguments {
				helpMessage += fmt.Sprintf("    %s: %s\n", arg.Name, arg.Description)
			}
		}

		helpMessage += "\n"
	}

	return helpMessage
}
