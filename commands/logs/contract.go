package logs

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/commands/sites"
	"github.com/daticahealth/cli/config"
	"github.com/daticahealth/cli/lib/auth"
	"github.com/daticahealth/cli/lib/jobs"
	"github.com/daticahealth/cli/lib/prompts"
	"github.com/daticahealth/cli/models"
	"github.com/jault3/mow.cli"
)

type CMDLogQuery struct {
	Query   string // default *
	Follow  bool
	Hours   int
	Minutes int
	Seconds int
	Service string
	JobID   string
	Target  string
}

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "logs",
	ShortHelp: "Show the logs in your terminal streamed from your logging dashboard",
	LongHelp: "`logs` prints out your application logs directly from your logging Dashboard. " +
		"If you do not see your logs, try adjusting the number of hours, minutes, or seconds of logs that are retrieved with the `--hours`, `--minutes`, and `--seconds` options respectively. " +
		"To specify a specific service, job, or target use the '--service', '--job-id', and '--target' commands. " +
		"You must specify a service to use '--job-id' or '--target', and you cannot specify both a job-id and a target at the same time. " +
		"You can also follow the logs with the `-f` option. " +
		"When using `-f` all logs will be printed to the console within the given time frame as well as any new logs that are sent to the logging Dashboard for the duration of the command. " +
		"When using the `-f` option, hit ctrl-c to stop. Here are some sample commands\n\n" +
		"```\ndatica -E \"<your_env_name>\" logs --hours=6 --minutes=30\n" +
		"datica -E \"<your_env_name>\" logs -f\n" +
		"datica -E \"<your_env_name>\" logs --service=\"<your_service_name>\"\n" +
		"datica -E \"<your_env_name>\" logs --service=\"<your_service_name>\" --job-id=\"<your_job_id>\"\n```",
	// TODO: add documentation here
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			query := cmd.StringArg("QUERY", "*", "The query to send to your logging dashboard's elastic search (regex is supported)")
			follow := cmd.BoolOpt("f follow", false, "Tail/follow the logs (Equivalent to -t)")
			tail := cmd.BoolOpt("t tail", false, "Tail/follow the logs (Equivalent to -f)")
			hours := cmd.IntOpt("hours", 0, "The number of hours before now (in combination with minutes and seconds) to retrieve logs")
			mins := cmd.IntOpt("minutes", 0, "The number of minutes before now (in combination with hours and seconds) to retrieve logs")
			secs := cmd.IntOpt("seconds", 0, "The number of seconds before now (in combination with hours and minutes) to retrieve logs")
			service := cmd.StringOpt("service", "", "Query logs for a specific service label")
			jobID := cmd.StringOpt("job-id", "", "Query logs for a particular job by id")
			target := cmd.StringOpt("target", "", "Query logs for a particular procfile target")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(settings); err != nil {
					logrus.Fatal(err.Error())
				}
				cmdQuery := CMDLogQuery{
					Query:   *query,
					Follow:  *follow || *tail,
					Hours:   *hours,
					Minutes: *mins,
					Seconds: *secs,
					Service: *service,
					JobID:   *jobID,
					Target:  *target,
				}
				err := CmdLogs(&cmdQuery, settings.EnvironmentID, settings, New(settings), prompts.New(), environments.New(settings), services.New(settings), jobs.New(settings), sites.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "[QUERY] [(-f | -t)] [--hours] [--minutes] [--seconds] [--service [(--job-id | --target)]]"
		}
	},
}

type queryGenerator func(queryString, appLogsIdentifier, appLogsValue string, timestamp time.Time, from int, hostNames []string, fileName string) ([]byte, error)

// ILogs ...
type ILogs interface {
	Output(queryString, domain string, generator queryGenerator, from int, startTimestamp time.Time, endTimestamp time.Time, hostNames []string, fileName string) (int, error)
	RetrieveElasticsearchVersion(domain string) (string, error)
	Stream(queryString, domain string, generator queryGenerator, from int, timestamp time.Time, hostNames []string, fileName string) error
	Watch(queryString, domain string) error
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
