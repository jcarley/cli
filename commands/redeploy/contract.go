package redeploy

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "redeploy",
	ShortHelp: "Redeploy a service without having to do a git push",
	LongHelp: "`redeploy` deploys an identical copy of the given service. " +
		"For code services, this avoids having to perform a code push. You skip the git push and the build. " +
		"For service proxies, new instances simply replace the old ones. " +
		"All other service types cannot be redeployed with this command. Here is a sample command\n\n" +
		"```\ncatalyze redeploy app01\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to redeploy (i.e. 'app01')")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRedeploy(settings.EnvironmentID, *serviceName, jobs.New(settings), services.New(settings), environments.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "SERVICE_NAME"
		}
	},
}
