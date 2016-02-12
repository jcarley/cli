package metrics

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
	Name:      "metrics",
	ShortHelp: "Print service and environment metrics in your local time zone",
	LongHelp:  "Print service and environment metrics in your local time zone",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to print metrics for")
			json := cmd.BoolOpt("json", false, "Output the data as json")
			csv := cmd.BoolOpt("csv", false, "Output the data as csv")
			spark := cmd.BoolOpt("spark", false, "Output the data using spark lines")
			stream := cmd.BoolOpt("stream", false, "Repeat calls once per minute until this process is interrupted.")
			mins := cmd.IntOpt("m mins", 1, "How many minutes worth of logs to retrieve.")
			cmd.Action = func() {
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Println(err.Error())
					os.Exit(1)
				}
				/*err := CmdMetrics(*serviceName, *json, *csv, *spark, *stream, *mins, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}*/
				logrus.Debugf("%s, %t, %t, %t, %t, %d", serviceName, json, csv, spark, stream, mins)
				logrus.Fatal("shits broke yo")
			}
			cmd.Spec = "[SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]"
		}
	},
}

// IMetrics
type IMetrics interface {
	// TODO rework metrics and remove this generic metrics command. most likely
	// the Text, CSV, JSON, and Spark will turn into their own Cmd methods.
	// such as CmdText and CmdCSV. Then just API calls go in this interface.
	Metrics(jsonFlag bool, csvFlag bool, sparkFlag bool, streamFlag bool, mins int, service *models.Service) error
	Text() error
	CSV() error
	JSON() error
	Spark() error
	RetrieveEnvironmentMetrics(mins int) (*[]models.Metrics, error)
	RetrieveServiceMetrics(mins int) (*models.Metrics, error)
}

// SMetrics is a concrete implementation of IMetrics
type SMetrics struct {
	Settings *models.Settings
}

// New returns an instance of IMetrics
func New(settings *models.Settings) IMetrics {
	return &SMetrics{
		Settings: settings,
	}
}
