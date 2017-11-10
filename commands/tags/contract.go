package tags

import (
	"github.com/Sirupsen/logrus"
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
	LongHelp: "`tags` allows interactions with tags. " +
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
		"```\ndatica -E \"<your_env_name>\" tags list pod012345/my-image\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			image := cmd.StringArg("IMAGE_NAME", "", "The name of the image to list tags for, including the environment's namespace.")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := cmdTagList(images.New(settings), *image)
				if err != nil {
					logrus.Fatalln(err.Error())
				}
			}
		}
	},
}

var deleteCmd = models.Command{
	Name:      "delete",
	ShortHelp: "Delete a tag for a given image",
	LongHelp: "Delete a tag for a given image. Example:\n" +
		"```\ndatica -E \"<your_env_name>\" tags delete pod012345/my-image v1\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			image := cmd.StringArg("IMAGE_NAME", "", "The name of the image to list tags for, including the environment's namespace.")
			tag := cmd.StringArg("TAG", "", "The tag to delete.")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := cmdTagDelete(images.New(settings), *image, *tag)
				if err != nil {
					logrus.Fatalln(err.Error())
				}
			}
		}
	},
}
