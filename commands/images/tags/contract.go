package tags

import (
	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/config"
	"github.com/daticahealth/cli/lib/auth"
	"github.com/daticahealth/cli/lib/images"
	"github.com/daticahealth/cli/lib/prompts"
	"github.com/daticahealth/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "tags",
	ShortHelp: "Operations for working with tags",
	LongHelp: "<code>tags</code> allows interactions with container image tags. " +
		"This command cannot be run directly, but has subcommands.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.CommandLong(listCmd.Name, listCmd.ShortHelp, listCmd.LongHelp, listCmd.CmdFunc(settings))
			cmd.CommandLong(deleteCmd.Name, deleteCmd.ShortHelp, deleteCmd.LongHelp, deleteCmd.CmdFunc(settings))
		}
	},
}

var listCmd = models.Command{
	Name:      "list",
	ShortHelp: "List tags for a given image",
	LongHelp: "List pushed tags for given image. Example:\n" +
		"<pre>\ndatica -E \"<your_env_name>\" images tags list <image>\n</pre>",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			image := cmd.StringArg("IMAGE_NAME", "", "The name of the image to list tags for. (e.g. 'my-image')")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := cmdTagList(images.New(settings), environments.New(settings), settings.EnvironmentID, *image)
				if err != nil {
					logrus.Fatalln(err.Error())
				}
			}
		}
	},
}

var deleteCmd = models.Command{
	Name:      "rm",
	ShortHelp: "Delete a tag for a given image",
	LongHelp: "Delete a tag for a given image. Example:\n" +
		"<pre>\ndatica -E \"<your_env_name>\" images tags rm <image>:<tag>\n</pre>",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			image := cmd.StringArg("TAGGED_IMAGE", "", "The name and tag of the image to delete. (e.g. 'my-image:tag')")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := cmdTagDelete(images.New(settings), prompts.New(), environments.New(settings), settings.EnvironmentID, *image)
				if err != nil {
					logrus.Fatalln(err.Error())
				}
			}
		}
	},
}
