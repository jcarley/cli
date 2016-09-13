package worker

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
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
	Name:      "worker",
	ShortHelp: "Manage a service's workers",
	LongHelp:  "Manage a service's workers",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(DeploySubCmd.Name, DeploySubCmd.ShortHelp, DeploySubCmd.CmdFunc(settings))
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(RmSubCmd.Name, RmSubCmd.ShortHelp, RmSubCmd.CmdFunc(settings))
			cmd.Command(ScaleSubCmd.Name, ScaleSubCmd.ShortHelp, ScaleSubCmd.CmdFunc(settings))

			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to use to start a worker. Defaults to the associated service.")
			target := cmd.StringArg("TARGET", "", "The name of the Procfile target to invoke as a worker")
			cmd.Action = func() {
				logrus.Warnln("This command has been moved! Please use \"catalyze worker deploy\" instead. This alias will be removed in the next CLI update.")
				logrus.Warnln("You can list all available worker subcommands by running \"catalyze worker --help\".")
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				if *target == "" {
					logrus.Fatal("TARGET is a required argument")
				}
				err := CmdWorker(*serviceName, settings.ServiceID, *target, New(settings), services.New(settings), jobs.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "[SERVICE_NAME] [TARGET]"
		}
	},
}

var DeploySubCmd = models.Command{
	Name:      "deploy",
	ShortHelp: "Deploy new workers for a given service",
	LongHelp:  "Deploy new workers for a given service",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to use to deploy a worker")
			target := subCmd.StringArg("TARGET", "", "The name of the Procfile target to invoke as a worker")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdDeploy(*serviceName, *target, New(settings), services.New(settings), jobs.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "SERVICE_NAME TARGET"
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "Lists all workers for a given service",
	LongHelp:  "Lists all workers for a given service",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to list workers for")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdList(*serviceName, New(settings), services.New(settings), jobs.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "SERVICE_NAME"
		}
	},
}

var RmSubCmd = models.Command{
	Name:      "rm",
	ShortHelp: "Remove all workers for a given service and target",
	LongHelp:  "Remove all workers for a given service and target",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service running the workers")
			target := subCmd.StringArg("TARGET", "", "The worker target to remove")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRm(*serviceName, *target, New(settings), services.New(settings), prompts.New(), jobs.New(settings))
				if err != nil {
					logrus.Fatalln(err.Error())
				}
			}
			subCmd.Spec = "SERVICE_NAME TARGET"
		}
	},
}

var ScaleSubCmd = models.Command{
	Name:      "scale",
	ShortHelp: "Scale existing workers up or down for a given service and target",
	LongHelp:  "Scale existing workers up or down for a given service and target",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service running the workers")
			target := subCmd.StringArg("TARGET", "", "The worker target to scale up or down")
			scale := subCmd.StringArg("SCALE", "", "The new scale (or change in scale) for the given worker target. This can be a single value (i.e. 2) representing the final number of workers that should be running. Or this can be a change represented by a plus or minus sign followed by the changed (i.e. +2 or -1)")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdScale(*serviceName, *target, *scale, New(settings), services.New(settings), prompts.New(), jobs.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "SERVICE_NAME TARGET SCALE"
		}
	},
}

// IWorker
type IWorker interface {
	ParseScale(scaleString string) (func(scale, change int) int, int, error)
	Retrieve(svcID string) (*models.Workers, error)
	Update(svcID string, workers *models.Workers) error
}

// SWorker is a concrete implementation of IWorker
type SWorker struct {
	Settings *models.Settings
}

// New returns an instance of IWorker
func New(settings *models.Settings) IWorker {
	return &SWorker{
		Settings: settings,
	}
}
