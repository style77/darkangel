package core

import (
	"darkangel/server/constant"
	"fmt"
)

func HandleCommandDoesNotExists(commandName string) CommandResult {
	return CommandResult{
		CommandName: commandName,
		Result:      false,
		Output:      fmt.Sprintf(constant.Red+"Command %s doesn't exists."+constant.Reset, commandName),
	}
}

func HandleCommandWrongArgumentsLength(command Command, request Request) CommandResult {
	return CommandResult{
		CommandName: command.Name,
		Result:      false,
		Output:      fmt.Sprintf(constant.Red+"Command %s takes %d argument. Passed %d."+constant.Reset, command.Name, len(command.Arguments), len(request.Args)),
	}
}

func HandleWrongTarget(request Request) CommandResult {
	return CommandResult{
		CommandName: request.CommandName,
		Result:      false,
		Output:      fmt.Sprintf(constant.Red+"Target \"%s\" is either offline or wrong."+constant.Reset, request.Args[0]),
	}
}
