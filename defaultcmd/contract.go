package defaultcmd

import (
	"fmt"
	"os"

	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "default",
	ShortHelp: "Set the default associated environment",
	LongHelp:  "Set the default associated environment",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			envAlias := cmd.StringArg("ENV_ALIAS", "", "The alias of an already associated environment to set as the default")
			cmd.Action = func() {
				id := New(settings, *envAlias)
				err := id.Set()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			cmd.Spec = "ENV_ALIAS"
		}
	},
}

// IDefault
type IDefault interface {
	Set() error
}

// SDefault is a concrete implementation of IDefault
type SDefault struct {
	Settings *models.Settings

	Alias string
}

// New returns an instance of IDefault
func New(settings *models.Settings, alias string) IDefault {
	return &SDefault{
		Settings: settings,
		Alias:    alias,
	}
}
