package logout

import (
	"fmt"
	"os"

	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "logout",
	ShortHelp: "Clear the stored user information from your local machine",
	LongHelp:  "Clear the stored user information from your local machine",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				err := CmdLogout(New(settings))
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
		}
	},
}

// ILogout
type ILogout interface {
	Logout() error
}

// SLogout is a concrete implementation of ILogout
type SLogout struct {
	Settings *models.Settings
}

// New returns an instance of ILogout
func New(settings *models.Settings) ILogout {
	return &SLogout{
		Settings: settings,
	}
}
