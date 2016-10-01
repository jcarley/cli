package domain

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/commands/sites"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "domain",
	ShortHelp: "Print out the temporary domain name of the environment",
	LongHelp: "`domain` prints out the temporary domain name setup by Catalyze for an environment. " +
		"This domain name typically takes the form podXXXXX.catalyzeapps.com but may vary based on the environment. Here is a sample command\n\n" +
		"```\ncatalyze domain\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdDomain(settings.EnvironmentID, environments.New(settings), services.New(settings), sites.New(settings))
				if err != nil {
					logrus.Fatalln(err.Error())
				}
			}
		}
	},
}
