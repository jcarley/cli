package vars

import (
	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/config"
	"github.com/daticahealth/cli/lib/auth"
	"github.com/daticahealth/cli/lib/prompts"
	"github.com/daticahealth/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "vars",
	ShortHelp: "Interaction with environment variables for an environment",
	LongHelp:  "The <code>vars</code> command allows you to manage environment variables for your code services. The vars command can not be run directly but has subcommands.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.CommandLong(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.LongHelp, ListSubCmd.CmdFunc(settings))
			cmd.CommandLong(SetSubCmd.Name, SetSubCmd.ShortHelp, SetSubCmd.LongHelp, SetSubCmd.CmdFunc(settings))
			cmd.CommandLong(UnsetSubCmd.Name, UnsetSubCmd.ShortHelp, UnsetSubCmd.LongHelp, UnsetSubCmd.CmdFunc(settings))
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all environment variables",
	LongHelp: "<code>vars list</code> prints out all known environment variables for the given code service. " +
		"You can print out environment variables in JSON or YAML format through the <code>--json</code> or <code>--yaml</code> flags. " +
		"Here are some sample commands\n\n" +
		"<pre>\ndatica -E \"<your_env_name>\" vars list code-1\n" +
		"datica -E \"<your_env_name>\" vars list code-1 --json\n</pre>",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service containing the environment variables.")
			json := subCmd.BoolOpt("json", false, "Output environment variables in JSON format")
			yaml := subCmd.BoolOpt("yaml", false, "Output environment variables in YAML format")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(settings); err != nil {
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
				err := CmdList(*serviceName, formatter, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "SERVICE_NAME [--json | --yaml]"
		}
	},
}

var SetSubCmd = models.Command{
	Name:      "set",
	ShortHelp: "Set one or more new environment variables or update the values of existing ones",
	LongHelp: "<code>vars set</code> allows you to add new environment variables or update the value of an existing environment variable on the given code service. " +
		"You can set/update 1 or more environment variables at a time with this command by repeating the <code>-v</code> option multiple times. " +
		"Once new environment variables are added or values updated, a redeploy is required for the given code service to have access to the new values. " +
		"The environment variables must be of the form <code><key>=<value></code>. Here is a sample command\n\n" +
		"<pre>\ndatica -E \"<your_env_name>\" vars set code-1 -v AWS_ACCESS_KEY_ID=1234 -v AWS_SECRET_ACCESS_KEY=5678\n</pre>",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service on which the environment variables will be set.")
			variables := subCmd.Strings(cli.StringsOpt{
				Name:      "v variable",
				Value:     []string{},
				Desc:      "The env variable to set or update in the form \"<key>=<value>\"",
				HideValue: true,
			})
			fileName := subCmd.StringOpt("f file", "", "The path to a file to import environment variables from. This file can be in JSON, YAML, or KEY=VALUE format")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdSet(*serviceName, *variables, *fileName, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "SERVICE_NAME (-v... | -f)"
		}
	},
}

var UnsetSubCmd = models.Command{
	Name:      "unset",
	ShortHelp: "Unset (delete) an existing environment variable",
	LongHelp: "<code>vars unset</code> removes environment variables from the given code service. " +
		"Only the environment variable name is required to unset. " +
		"Once environment variables are unset, a redeploy is required for the given code service to realize the variable was removed. " +
		"You can unset any number of environment variables in one command. " +
		"Here is a sample command\n\n" +
		"<pre>\ndatica -E \"<your_env_name>\" vars unset code-1 AWS_ACCESS_KEY_ID AWS_SECRET_ACCES_KEY_ID\n</pre>",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service on which the environment variables will be unset.")
			variables := subCmd.Strings(cli.StringsArg{
				Name:      "VARIABLE",
				Value:     []string{},
				Desc:      "The names of environment variables to unset",
				HideValue: true,
			})
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdUnset(*serviceName, *variables, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "SERVICE_NAME VARIABLE..."
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
