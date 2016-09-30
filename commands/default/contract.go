package defaultcmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "default",
	ShortHelp: "[DEPRECATED] Set the default associated environment",
	LongHelp: "The `default` command has been deprecated! It will be removed in a future version. Please specify `-E` on all commands instead of using the default.\n\n" +
		"`default` sets the default environment for all commands that don't specify an environment with the `-E` flag. " +
		"See [scope](#global-scope) for more information on scope and default environments. " +
		"When setting a default environment, you must give the alias of the environment if one was set when it was associated and not the real environment name. Here is a sample command\n\n" +
		"```catalyze default prod```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			alias := cmd.StringArg("ENV_ALIAS", "", "The alias of an already associated environment to set as the default")
			cmd.Action = func() {
				if err := config.CheckRequiredAssociation(true, false, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdDefault(*alias, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "ENV_ALIAS"
		}
	},
}

// IDefault
type IDefault interface {
	Set(alias string) error
}

// SDefault is a concrete implementation of IDefault
type SDefault struct {
	Settings *models.Settings
}

// New returns an instance of IDefault
func New(settings *models.Settings) IDefault {
	return &SDefault{
		Settings: settings,
	}
}
