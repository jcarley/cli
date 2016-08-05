package environments

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "environments",
	ShortHelp: "Manage environments for which you have access",
	LongHelp:  "Manage environments for which you have access",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(RenameSubCmd.Name, RenameSubCmd.ShortHelp, RenameSubCmd.CmdFunc(settings))
			cmd.Action = func() {
				logrus.Warnln("This command has been moved! Please use \"catalyze environments list\" instead. This alias will be removed in the next CLI update.")
				logrus.Warnln("You can list all available environments subcommands by running \"catalyze environments --help\".")
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdList(New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all environments you have access to",
	LongHelp:  "List all environments you have access to",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatalln(err.Error())
				}
				err := CmdList(New(settings))
				if err != nil {
					logrus.Fatalln(err.Error())
				}
			}
		}
	},
}

var RenameSubCmd = models.Command{
	Name:      "rename",
	ShortHelp: "Rename an environment",
	LongHelp:  "Rename an environment",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			name := subCmd.StringArg("NAME", "", "The new name of the environment")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRename(settings.EnvironmentID, *name, New(settings))
				if err != nil {
					logrus.Fatalln(err.Error())
				}
			}
			subCmd.Spec = "NAME"
		}
	},
}

// IEnvironments is an interface for interacting with environments
type IEnvironments interface {
	List() (*[]models.Environment, error)
	Retrieve(envID string) (*models.Environment, error)
	Update(envID string, updates map[string]string) error
}

// SEnvironments is a concrete implementation of IEnvironments
type SEnvironments struct {
	Settings *models.Settings
}

// New generates a new instance of IEnvironments
func New(settings *models.Settings) IEnvironments {
	return &SEnvironments{
		Settings: settings,
	}
}
