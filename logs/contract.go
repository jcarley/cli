package logs

import (
	"fmt"
	"os"

	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "logs",
	ShortHelp: "Show the logs in your terminal streamed from your logging dashboard",
	LongHelp:  "Show the logs in your terminal streamed from your logging dashboard",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			query := cmd.StringArg("QUERY", "*", "The query to send to your logging dashboard's elastic search (regex is supported)")
			follow := cmd.BoolOpt("f follow", false, "Tail/follow the logs (Equivalent to -t)")
			tail := cmd.BoolOpt("t tail", false, "Tail/follow the logs (Equivalent to -f)")
			hours := cmd.IntOpt("hours", 0, "The number of hours before now (in combination with minutes and seconds) to retrieve logs")
			mins := cmd.IntOpt("minutes", 1, "The number of minutes before now (in combination with hours and seconds) to retrieve logs")
			secs := cmd.IntOpt("seconds", 0, "The number of seconds before now (in combination with hours and minutes) to retrieve logs")
			cmd.Action = func() {
				err := CmdLogs(*query, *follow || *tail, *hours, *mins, *secs, New(settings))
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			cmd.Spec = "[QUERY] [(-f | -t)] [--hours] [--minutes] [--seconds]"
		}
	},
}

// ILogs
type ILogs interface {
	Output(queryString string, tail bool, hours, minutes, seconds int) error
	Stream() error
}

// SLogs is a concrete implementation of ILogs
type SLogs struct {
	Settings *models.Settings
}

// New returns an instance of ILogs
func New(settings *models.Settings) ILogs {
	return &SLogs{
		Settings: settings,
	}
}
