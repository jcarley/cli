package catalyze

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/catalyzeio/cli/commands/associate"
	"github.com/catalyzeio/cli/commands/associated"
	"github.com/catalyzeio/cli/commands/console"
	"github.com/catalyzeio/cli/commands/dashboard"
	"github.com/catalyzeio/cli/commands/db"
	"github.com/catalyzeio/cli/commands/default"
	"github.com/catalyzeio/cli/commands/disassociate"
	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/commands/files"
	"github.com/catalyzeio/cli/commands/invites"
	"github.com/catalyzeio/cli/commands/keys"
	"github.com/catalyzeio/cli/commands/logout"
	"github.com/catalyzeio/cli/commands/logs"
	"github.com/catalyzeio/cli/commands/metrics"
	"github.com/catalyzeio/cli/commands/rake"
	"github.com/catalyzeio/cli/commands/redeploy"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/commands/ssl"
	"github.com/catalyzeio/cli/commands/status"
	"github.com/catalyzeio/cli/commands/supportids"
	"github.com/catalyzeio/cli/commands/update"
	"github.com/catalyzeio/cli/commands/users"
	"github.com/catalyzeio/cli/commands/vars"
	"github.com/catalyzeio/cli/commands/whoami"
	"github.com/catalyzeio/cli/commands/worker"

	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/models"

	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/pods"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/lib/updater"

	"github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
)

type simpleLogger struct{}

func (s *simpleLogger) Format(entry *logrus.Entry) ([]byte, error) {
	//	entry.Data["msg"]`. The message passed from Info, Warn, Error ..
	// * `entry.Data["time"]`. The timestamp.
	// * `entry.Data["level"]. The level the entry was logged at.
	levelString := fmt.Sprintf("[%s] ", entry.Level)
	if entry.Level == logrus.InfoLevel {
		levelString = ""
	}
	l := fmt.Sprintf("%s%s\n", levelString, entry.Message)
	return []byte(l), nil
}

// Run runs the Catalyze CLI
func Run() {
	if updater.AutoUpdater != nil {
		updater.AutoUpdater.BackgroundRun()
	}

	logrus.SetFormatter(&simpleLogger{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(config.LogLevel)

	var app = cli.App("catalyze", fmt.Sprintf("Catalyze CLI. Version %s", config.VERSION))

	authHost := os.Getenv(config.AuthHostEnvVar)
	if authHost == "" {
		authHost = config.AuthHost
	}
	paasHost := os.Getenv(config.PaasHostEnvVar)
	if paasHost == "" {
		paasHost = config.PaasHost
	}
	username := app.String(cli.StringOpt{
		Name:      "U username",
		Desc:      "Catalyze Username",
		EnvVar:    config.CatalyzeUsernameEnvVar,
		HideValue: true,
	})
	password := app.String(cli.StringOpt{
		Name:      "P password",
		Desc:      "Catalyze Password",
		EnvVar:    config.CatalyzePasswordEnvVar,
		HideValue: true,
	})
	givenEnvName := app.String(cli.StringOpt{
		Name:      "E env",
		Desc:      "The local alias of the environment in which this command will be run",
		EnvVar:    config.CatalyzeEnvironmentEnvVar,
		HideValue: true,
	})
	if loggingLevel, err := strconv.ParseInt(os.Getenv(config.LogLevelEnvVar), 10, 64); err == nil {
		logrus.SetLevel(logrus.Level(loggingLevel))
	}
	settings := &models.Settings{}

	app.Before = func() {
		r := config.FileSettingsRetriever{}
		*settings = *r.GetSettings(*givenEnvName, "", authHost, paasHost, *username, *password)
		logrus.Debugf("%+v", settings)

		if settings.Pods == nil || len(*settings.Pods) == 0 {
			p := pods.New(settings)
			pods, err := p.List()
			if err == nil {
				settings.Pods = pods
				fmt.Println(settings.Pods)
			} else {
				logrus.Debugln(err.Error())
				// TODO check in the cmd wherever settings.Pods is used and check for null/empty
				// log this err
			}
		}

		a := auth.New(settings, prompts.New())
		user, err := a.Signin()
		if err != nil {
			logrus.Fatalln(err.Error())
		}
		settings.SessionToken = user.SessionToken
		settings.Username = user.Username
		settings.UsersID = user.UsersID
	}
	app.After = func() {
		config.SaveSettings(settings)
	}

	InitCLI(app, settings)

	archString := "other"
	switch runtime.GOARCH {
	case "386":
		archString = "32-bit"
	case "amd64":
		archString = "64-bit"
	case "arm":
		archString = "arm"
	}
	versionString := fmt.Sprintf("version %s %s", config.VERSION, archString)
	logrus.Debugln(versionString)
	app.Version("v version", versionString)
	app.Command("version", "Output the version and quit", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			fmt.Println(versionString)
		}
	})

	app.Run(os.Args)
}

// InitCLI adds arguments and commands to the given cli instance
func InitCLI(app *cli.Cli, settings *models.Settings) {
	app.Command(associate.Cmd.Name, associate.Cmd.ShortHelp, associate.Cmd.CmdFunc(settings))
	app.Command(associated.Cmd.Name, associated.Cmd.ShortHelp, associated.Cmd.CmdFunc(settings))
	app.Command(console.Cmd.Name, console.Cmd.ShortHelp, console.Cmd.CmdFunc(settings))
	app.Command(dashboard.Cmd.Name, dashboard.Cmd.ShortHelp, dashboard.Cmd.CmdFunc(settings))
	app.Command(db.Cmd.Name, db.Cmd.ShortHelp, db.Cmd.CmdFunc(settings))
	app.Command(defaultcmd.Cmd.Name, defaultcmd.Cmd.ShortHelp, defaultcmd.Cmd.CmdFunc(settings))
	app.Command(disassociate.Cmd.Name, disassociate.Cmd.ShortHelp, disassociate.Cmd.CmdFunc(settings))
	app.Command(environments.Cmd.Name, environments.Cmd.ShortHelp, environments.Cmd.CmdFunc(settings))
	app.Command(files.Cmd.Name, files.Cmd.ShortHelp, files.Cmd.CmdFunc(settings))
	app.Command(invites.Cmd.Name, invites.Cmd.ShortHelp, invites.Cmd.CmdFunc(settings))
	app.Command(keys.Cmd.Name, keys.Cmd.ShortHelp, keys.Cmd.CmdFunc(settings))
	app.Command(logs.Cmd.Name, logs.Cmd.ShortHelp, logs.Cmd.CmdFunc(settings))
	app.Command(logout.Cmd.Name, logout.Cmd.ShortHelp, logout.Cmd.CmdFunc(settings))
	app.Command(metrics.Cmd.Name, metrics.Cmd.ShortHelp, metrics.Cmd.CmdFunc(settings))
	app.Command(rake.Cmd.Name, rake.Cmd.ShortHelp, rake.Cmd.CmdFunc(settings))
	app.Command(redeploy.Cmd.Name, redeploy.Cmd.ShortHelp, redeploy.Cmd.CmdFunc(settings))
	app.Command(services.Cmd.Name, services.Cmd.ShortHelp, services.Cmd.CmdFunc(settings))
	app.Command(ssl.Cmd.Name, ssl.Cmd.ShortHelp, ssl.Cmd.CmdFunc(settings))
	app.Command(status.Cmd.Name, status.Cmd.ShortHelp, status.Cmd.CmdFunc(settings))
	app.Command(supportids.Cmd.Name, supportids.Cmd.ShortHelp, supportids.Cmd.CmdFunc(settings))
	app.Command(update.Cmd.Name, update.Cmd.ShortHelp, update.Cmd.CmdFunc(settings))
	app.Command(users.Cmd.Name, users.Cmd.ShortHelp, users.Cmd.CmdFunc(settings))
	app.Command(vars.Cmd.Name, vars.Cmd.ShortHelp, vars.Cmd.CmdFunc(settings))
	app.Command(whoami.Cmd.Name, whoami.Cmd.ShortHelp, whoami.Cmd.CmdFunc(settings))
	app.Command(worker.Cmd.Name, worker.Cmd.ShortHelp, worker.Cmd.CmdFunc(settings))
}
