package deploy

import (
	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/config"
	"github.com/daticahealth/cli/lib/auth"
	"github.com/daticahealth/cli/lib/images"
	"github.com/daticahealth/cli/lib/jobs"
	"github.com/daticahealth/cli/lib/prompts"
	"github.com/daticahealth/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "deploy",
	ShortHelp: "Deploy a Docker image to a container service.",
	LongHelp: "<code>deploy</code> deploys a Docker image for the given service. " +
		"This command will only deploy for \"container\" services. " +
		"Here is a sample command\n\n" +
		"<pre>\ndatica -E \"<your_env_name>\" deploy <service> <image>:<tag>\n</pre>",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service where the image will be deployed. (e.g. 'container-1')")
			imageName := cmd.StringArg("TAGGED_IMAGE", "", "The name and tag of the image to deploy. (e.g. 'my-image:tag)")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdDeploy(settings.EnvironmentID, *serviceName, *imageName, jobs.New(settings), services.New(settings), environments.New(settings), images.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "SERVICE_NAME TAGGED_IMAGE"
		}
	},
}
