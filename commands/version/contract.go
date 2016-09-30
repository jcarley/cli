package version

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "version",
	ShortHelp: "Output the version and quit",
	LongHelp: "`version` prints out the current CLI version as well as the architecture it was built for (64-bit or 32-bit). " +
		"This is useful to see if you have the latest version of the CLI and when working with Catalyze support engineers to ensure you have the correct CLI installed. " +
		"Here is a sample command\n\n" +
		"```catalyze version```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				err := CmdVersion()
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}
