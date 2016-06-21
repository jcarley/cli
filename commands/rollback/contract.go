package rollback

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/redeploy"
	"github.com/catalyzeio/cli/commands/releases"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "rollback",
	ShortHelp: "Rollback a code service to a specific release",
	LongHelp:  "Rollback a code service to a specific release",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to rollback")
			releaseName := cmd.StringArg("RELEASE_NAME", "", "The name of the release to rollback to")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRollback(*serviceName, *releaseName, redeploy.New(settings), releases.New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "SERVICE_NAME RELEASE_NAME"
		}
	},
}
