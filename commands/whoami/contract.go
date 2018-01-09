package whoami

import (
	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/lib/auth"
	"github.com/daticahealth/cli/lib/prompts"
	"github.com/daticahealth/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "whoami",
	ShortHelp: "Retrieve your user ID",
	LongHelp: "<code>whoami</code> prints out the currently logged in user's users ID. " +
		"This is used with Datica support engineers. Here is a sample command\n\n" +
		"<pre>\ndatica whoami\n</pre>",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdWhoAmI(New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

// IWhoAmI
type IWhoAmI interface {
	WhoAmI() (string, error)
}

// SWhoAmI is a concrete implementation of IWhoAmI
type SWhoAmI struct {
	Settings *models.Settings
}

// New returns an instance of IWhoAmI
func New(settings *models.Settings) IWhoAmI {
	return &SWhoAmI{
		Settings: settings,
	}
}
