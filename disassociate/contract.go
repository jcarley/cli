package disassociate

import (
	"fmt"
	"os"

	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "disassociate",
	ShortHelp: "Remove the association with an environment",
	LongHelp:  "Remove the association with an environment",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			envAlias := cmd.StringArg("ENV_ALIAS", "", "The alias of an already associated environment to disassociate")
			cmd.Action = func() {
				id := New(settings, *envAlias)
				err := id.Disassociate()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			cmd.Spec = "ENV_ALIAS"
		}
	},
}

// IDisassociate
type IDisassociate interface {
	Disassociate() error
}

// SDisassociate is a concrete implementation of IDisassociate
type SDisassociate struct {
	Settings *models.Settings

	Alias string
}

// New returns an instance of IDisassociate
func New(settings *models.Settings, alias string) IDisassociate {
	return &SDisassociate{
		Settings: settings,
		Alias:    alias,
	}
}
