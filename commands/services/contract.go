package services

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "services",
	ShortHelp: "Perform operations on an environment's services",
	LongHelp:  "Perform operations on an environment's services",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(StopSubCmd.Name, StopSubCmd.ShortHelp, StopSubCmd.CmdFunc(settings))
			cmd.Action = func() {
				logrus.Warnln("This command has been moved! Please use \"catalyze services list\" instead. This alias will be removed in the next CLI update.")
				logrus.Warnln("You can list all available services subcommands by running \"catalyze services --help\".")
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdServices(New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all services for your environment",
	LongHelp:  "List all services for your environment",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdServices(New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

var StopSubCmd = models.Command{
	Name:      "stop",
	ShortHelp: "Stop all instances of a given service (including all workers, rake tasks, and open consoles)",
	LongHelp:  "Stop all instances of a given service (including all workers, rake tasks, and open consoles)",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			svcName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to stop")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdStop(*svcName, New(settings), jobs.New(settings), prompts.New())
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "SERVICE_NAME"
		}
	},
}

// IServices
type IServices interface {
	List() (*[]models.Service, error)
	ListByEnvID(envID, podID string) (*[]models.Service, error)
	Retrieve(svcID string) (*models.Service, error)
	RetrieveByLabel(label string) (*models.Service, error)
}

// SServices is a concrete implementation of IServices
type SServices struct {
	Settings *models.Settings
}

// New generates a new instance of IServices
func New(settings *models.Settings) IServices {
	return &SServices{
		Settings: settings,
	}
}
