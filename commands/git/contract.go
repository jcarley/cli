package git

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "git-remote",
	ShortHelp: "Manage git remotes to Catalyze code services",
	LongHelp: "The `git-remote` command allows you to interact with code service remote git URLs. " +
		"The git-remote command can not be run directly but has sub commands.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.CommandLong(AddSubCmd.Name, AddSubCmd.ShortHelp, AddSubCmd.LongHelp, AddSubCmd.CmdFunc(settings))
			cmd.CommandLong(ShowSubCmd.Name, ShowSubCmd.ShortHelp, ShowSubCmd.LongHelp, ShowSubCmd.CmdFunc(settings))
		}
	},
}

var AddSubCmd = models.Command{
	Name:      "add",
	ShortHelp: "Add the git remote for the given code service to the local git repo",
	LongHelp: "`git-remote add` adds the proper git remote to a local git repository with the given remote name and service. " +
		"Here is a sample command\n\n" +
		"```\ncatalyze git-remote add code-1 -r catalyze-code-1\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to add a git remote for")
			remote := subCmd.StringOpt("r remote", "catalyze", "The name of the git remote to be added")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdAdd(*serviceName, *remote, New(), services.New(settings))
				if err != nil {
					logrus.Fatalln(err.Error())
				}
			}
			subCmd.Spec = "SERVICE_NAME [-r]"
		}
	},
}

var ShowSubCmd = models.Command{
	Name:      "show",
	ShortHelp: "Print out the git remote for a given code service",
	LongHelp: "`git-remote show` prints out the git remote URL for the given service. " +
		"This can be used to do a manual push or use the git remote for another purpose such as a CI integration. " +
		"Here is a sample command\n\n" +
		"```\ncatalyze git-remote show code-1\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to add a git remote for")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdShow(*serviceName, services.New(settings))
				if err != nil {
					logrus.Fatalln(err.Error())
				}
			}
			subCmd.Spec = "SERVICE_NAME"
		}
	},
}

// IGit is an interface through which you can perform git operations
type IGit interface {
	Add(remote, gitURL string) error
	Exists() bool
	List() ([]string, error)
	Rm(remote string) error
}

// SGit is an implementor of IGit
type SGit struct{}

// New creates a new instance of IGit
func New() IGit {
	return &SGit{}
}
