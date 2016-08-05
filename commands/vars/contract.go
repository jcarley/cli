package vars

import (
	"github.com/Sirupsen/logrus"
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
	Name:      "vars",
	ShortHelp: "Interaction with environment variables for the associated environment",
	LongHelp:  "Interaction with environment variables for the associated environment",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(SetSubCmd.Name, SetSubCmd.ShortHelp, SetSubCmd.CmdFunc(settings))
			cmd.Command(UnsetSubCmd.Name, UnsetSubCmd.ShortHelp, UnsetSubCmd.CmdFunc(settings))
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all environment variables",
	LongHelp:  "List all environment variables",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service containing the environment variables. Defaults to the associated service.")
			json := subCmd.BoolOpt("json", false, "Output environment variables in JSON format")
			yaml := subCmd.BoolOpt("yaml", false, "Output environment variables in YAML format")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				var formatter Formatter
				if *json {
					formatter = &JSONFormatter{}
				} else if *yaml {
					formatter = &YAMLFormatter{}
				} else {
					formatter = &PlainFormatter{}
				}
				err := CmdList(*serviceName, settings.ServiceID, formatter, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "[SERVICE_NAME] [--json | --yaml]"
		}
	},
}

var SetSubCmd = models.Command{
	Name:      "set",
	ShortHelp: "Set one or more new environment variables or update the values of existing ones",
	LongHelp:  "Set one or more new environment variables or update the values of existing ones",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service on which the environment variables will be set. Defaults to the associated service.")
			variables := subCmd.Strings(cli.StringsOpt{
				Name:      "v variable",
				Value:     []string{},
				Desc:      "The env variable to set or update in the form \"<key>=<value>\"",
				HideValue: true,
			})
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdSet(*serviceName, settings.ServiceID, *variables, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "[SERVICE_NAME] -v..."
		}
	},
}

var UnsetSubCmd = models.Command{
	Name:      "unset",
	ShortHelp: "Unset (delete) an existing environment variable",
	LongHelp:  "Unset (delete) an existing environment variable",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service on which the environment variables will be unset. Defaults to the associated service.")
			variable := subCmd.StringArg("VARIABLE", "", "The name of the environment variable to unset")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdUnset(*serviceName, settings.ServiceID, *variable, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "[SERVICE_NAME] VARIABLE"
		}
	},
}

// IVars
type IVars interface {
	List(svcID string) (map[string]string, error)
	Set(svcID string, envVarsMap map[string]string) error
	Unset(svcID, key string) error
}

// SVars is a concrete implementation of IVars
type SVars struct {
	Settings *models.Settings
}

// New generates a new instance of IVars
func New(settings *models.Settings) IVars {
	return &SVars{
		Settings: settings,
	}
}
