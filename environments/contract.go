package environments

import (
	"fmt"
	"os"

	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "environments",
	ShortHelp: "List all environments you have access to",
	LongHelp:  "List all environments you have access to",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				ie := New(settings, "")
				err := CmdEnvironments(ie)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
		}
	},
}

// IEnvironments is an interface for interacting with environments
type IEnvironments interface {
	List() (*[]models.Environment, error)
	Retrieve() (*models.Environment, error)
}

// SEnvironments is a concrete implementation of IEnvironments
type SEnvironments struct {
	Settings *models.Settings

	EnvID string
}

// New generates a new instance of IEnvironments
func New(settings *models.Settings, envID string) IEnvironments {
	return &SEnvironments{
		Settings: settings,

		EnvID: envID,
	}
}
