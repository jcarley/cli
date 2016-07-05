package clear

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "clear",
	ShortHelp: "Clear out information in the global settings file to fix a misconfigured CLI.",
	LongHelp:  "Clear out information in the global settings file to fix a misconfigured CLI.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			privateKey := cmd.BoolOpt("private-key", false, "Clear out the saved private key information")
			session := cmd.BoolOpt("session", false, "Clear out all session information")
			envs := cmd.BoolOpt("environments", false, "Clear out all associated environments")
			defaultEnv := cmd.BoolOpt("default", false, "Clear out the saved default environment")
			pods := cmd.BoolOpt("pods", false, "Clear out all saved pods")
			all := cmd.BoolOpt("all", false, "Clear out all settings")
			cmd.Action = func() {
				if *all {
					*privateKey = true
					*session = true
					*envs = true
					*defaultEnv = true
					*pods = true
				}
				err := CmdClear(*privateKey, *session, *envs, *defaultEnv, *pods, settings)
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "[--private-key] [--session] [--environments] [--default] [--pods] [--all]"
		}
	},
}
