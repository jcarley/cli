package targets

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
	Name:      "targets",
	ShortHelp: "Operations for working with signed targets",
	LongHelp: "`targets` allows interactions with content verified targets in a repository. " +
		"This command cannot be run directly, but has subcommands.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.CommandLong(listCmd.Name, listCmd.ShortHelp, listCmd.LongHelp, listCmd.CmdFunc(settings))
			cmd.CommandLong(deleteCmd.Name, deleteCmd.ShortHelp, deleteCmd.LongHelp, deleteCmd.CmdFunc(settings))
			cmd.CommandLong(statusCmd.Name, statusCmd.ShortHelp, statusCmd.LongHelp, statusCmd.CmdFunc(settings))
			cmd.CommandLong(resetCmd.Name, resetCmd.ShortHelp, resetCmd.LongHelp, resetCmd.CmdFunc(settings))
		}
	},
}

var listCmd = models.Command{
	Name:      "list",
	ShortHelp: "List signed targets for an image",
	LongHelp: "`images targets list` lists signed targets for an image.\n" +
		"To search for a specific target, specify a tag with the image name in the format <image>:<tag>\n" +
		"Example:\n" +
		"```\ndatica -E \"<your_env_name>\" images targets list <image>\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			image := cmd.StringArg("IMAGE_NAME", "", "The name of the image to list targets for.")
			cmd.Action = func() {
				user, err := auth.New(settings, prompts.New()).Signin()
				if err != nil {
					logrus.Fatal(err.Error())
				}
				if err = config.CheckRequiredAssociation(settings); err != nil {
					logrus.Fatal(err.Error())
				}
				if err = cmdTargetsList(settings.EnvironmentID, *image, user, environments.New(settings), images.New(settings)); err != nil {
					logrus.Fatalln(err.Error())
				}
			}
		}
	},
}

var deleteCmd = models.Command{
	Name:      "rm",
	ShortHelp: "Delete a signed target for a given image",
	LongHelp: "Delete a signed target for a given image. You environment namespace will be filled in for you if not provided.\n" +
		"Example:\n" +
		"```\ndatica -E \"<your_env_name>\" images targets rm <image>:<tag>\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			image := cmd.StringArg("IMAGE_NAME", "", "The name of the image to delete targets for")
			cmd.Action = func() {
				user, err := auth.New(settings, prompts.New()).Signin()
				if err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(settings); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := cmdTargetsDelete(settings.EnvironmentID, *image, user, environments.New(settings), images.New(settings), prompts.New()); err != nil {
					logrus.Fatalln(err.Error())
				}
			}
		}
	},
}

var statusCmd = models.Command{
	Name:      "status",
	ShortHelp: "List local unpublished changes to the trust repository for an image",
	LongHelp: "`images targets status` lists unpublished changes to a local trust repository.\n" +
		"To search for changes to a specific target, specify a tag with the image name in the format <image>:<tag>\n" +
		"Example:\n" +
		"```\ndatica -E \"<your_env_name>\" images targets status <image>\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			image := cmd.StringArg("IMAGE_NAME", "", "The name of the image to list unpublished changes for.")
			cmd.Action = func() {
				user, err := auth.New(settings, prompts.New()).Signin()
				if err != nil {
					logrus.Fatal(err.Error())
				}
				if err = config.CheckRequiredAssociation(settings); err != nil {
					logrus.Fatal(err.Error())
				}
				if err = cmdTargetsStatus(settings.EnvironmentID, *image, user, environments.New(settings), images.New(settings)); err != nil {
					logrus.Fatalln(err.Error())
				}
			}
		}
	},
}

var resetCmd = models.Command{
	Name:      "reset",
	ShortHelp: "Clear unpublished changes in the local trust repository for an image",
	LongHelp: "`images targets reset` clears unpublished changes in a local trust repository. Does not affect your remote trust repository.\n" +
		"To reset changes for a specific target, specify a tag with the image name in the format <image>:<tag>\n" +
		"Example:\n" +
		"```\ndatica -E \"<your_env_name>\" images targets reset <image>\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			image := cmd.StringArg("IMAGE_NAME", "", "The name of the image to reset unpublished changes for.")
			cmd.Action = func() {
				user, err := auth.New(settings, prompts.New()).Signin()
				if err != nil {
					logrus.Fatal(err.Error())
				}
				if err = config.CheckRequiredAssociation(settings); err != nil {
					logrus.Fatal(err.Error())
				}
				if err = cmdTargetsReset(settings.EnvironmentID, *image, user, environments.New(settings), images.New(settings)); err != nil {
					logrus.Fatalln(err.Error())
				}
			}
		}
	},
}
