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
	ShortHelp: "Clear out information in the global settings file to fix a misconfigured CLI. All information will be cleared unless otherwise specified",
	LongHelp:  "Clear out information in the global settings file to fix a misconfigured CLI. All information will be cleared unless otherwise specified",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			privateKey := cmd.BoolOpt("private-key", true, "Clear out the saved private key information")
			session := cmd.BoolOpt("session", true, "Clear out all session information")
			envs := cmd.BoolOpt("environments", true, "Clear out all associated environments")
			defaultEnv := cmd.BoolOpt("default", true, "Clear out the saved default environment")
			pods := cmd.BoolOpt("pods", true, "Clear out all saved pods")
			cmd.Action = func() {
				err := CmdClear(*privateKey, *session, *envs, *defaultEnv, *pods, settings)
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "[--private-key] [--session] [--environments] [--default] [--pods]"
		}
	},
}
