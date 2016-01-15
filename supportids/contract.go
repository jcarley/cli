package supportids

import (
	"fmt"
	"os"

	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "supportids",
	ShortHelp: "Print out various IDs related to your associated environment to be used when contacting Catalyze support",
	LongHelp:  "Print out various IDs related to your associated environment to be used when contacting Catalyze support",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				err := CmdSupportIDs(New(settings))
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
		}
	},
}

// ISupportIDs
type ISupportIDs interface {
	SupportIDs() (string, string, string, error)
}

// SSupportIDs is a concrete implementation of ISupportIDs
type SSupportIDs struct {
	Settings *models.Settings
}

// New returns an instance of ISupportIDs
func New(settings *models.Settings) ISupportIDs {
	return &SSupportIDs{
		Settings: settings,
	}
}
