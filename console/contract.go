package console

import (
	"fmt"
	"os"

	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "console",
	ShortHelp: "Open a secure console to a service",
	LongHelp:  "Open a secure console to a service",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to open up a console for")
			command := cmd.StringArg("COMMAND", "", "An optional command to run when the console becomes available")
			cmd.Action = func() {
				ic := New(settings, *serviceName, *command)
				err := ic.Open()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			cmd.Spec = "SERVICE_NAME [COMMAND]"
		}
	},
}

// IConsole
type IConsole interface {
	Open() error
}

// SConsole is a concrete implementation of IConsole
type SConsole struct {
	Settings *models.Settings

	SvcName string
	Command string
}

// New returns an instance of IConsole
func New(settings *models.Settings, svcName, command string) IConsole {
	return &SConsole{
		Settings: settings,
		SvcName:  svcName,
		Command:  command,
	}
}
