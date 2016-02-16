package rake

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "rake",
	ShortHelp: "Execute a rake task",
	LongHelp:  "Execute a rake task",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			taskName := cmd.StringArg("TASK_NAME", "", "The name of the rake task to run")
			cmd.Action = func() {
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Println(err.Error())
					os.Exit(1)
				}
				err := CmdRake(*taskName, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "TASK_NAME"
		}
	},
}

// IRake
type IRake interface {
	Run(taskName string) error
}

// SRake is a concrete implementation of IRake
type SRake struct {
	Settings *models.Settings
}

// New returns an instance of IRake
func New(settings *models.Settings) IRake {
	return &SRake{
		Settings: settings,
	}
}
