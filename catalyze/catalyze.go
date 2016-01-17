package catalyze

import (
	"fmt"
	"os"
	"runtime"

	"github.com/catalyzeio/cli/associate"
	"github.com/catalyzeio/cli/associated"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/console"
	"github.com/catalyzeio/cli/dashboard"
	"github.com/catalyzeio/cli/db"
	"github.com/catalyzeio/cli/defaultcmd"
	"github.com/catalyzeio/cli/disassociate"
	"github.com/catalyzeio/cli/environments"
	"github.com/catalyzeio/cli/files"
	"github.com/catalyzeio/cli/invites"
	"github.com/catalyzeio/cli/logout"
	"github.com/catalyzeio/cli/logs"
	"github.com/catalyzeio/cli/metrics"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/rake"
	"github.com/catalyzeio/cli/redeploy"
	"github.com/catalyzeio/cli/services"
	"github.com/catalyzeio/cli/ssl"
	"github.com/catalyzeio/cli/status"
	"github.com/catalyzeio/cli/supportids"
	"github.com/catalyzeio/cli/update"
	"github.com/catalyzeio/cli/updater"
	"github.com/catalyzeio/cli/users"
	"github.com/catalyzeio/cli/vars"
	"github.com/catalyzeio/cli/whoami"
	"github.com/catalyzeio/cli/worker"
	"github.com/jawher/mow.cli"
)

// Run runs the Catalyze CLI
func Run() {
	if updater.AutoUpdater != nil {
		updater.AutoUpdater.BackgroundRun()
	}

	var app = cli.App("catalyze", fmt.Sprintf("Catalyze CLI. Version %s", config.VERSION))

	baasHost := os.Getenv("BAAS_HOST")
	if baasHost == "" {
		baasHost = config.BaasHost
	}
	paasHost := os.Getenv("PAAS_HOST")
	if paasHost == "" {
		paasHost = config.PaasHost
	}
	username := app.String(cli.StringOpt{
		Name:      "U username",
		Desc:      "Catalyze Username",
		EnvVar:    "CATALYZE_USERNAME",
		HideValue: true,
	})
	password := app.String(cli.StringOpt{
		Name:      "P password",
		Desc:      "Catalyze Password",
		EnvVar:    "CATALYZE_PASSWORD",
		HideValue: true,
	})
	givenEnvName := app.String(cli.StringOpt{
		Name:      "E env",
		Desc:      "The local alias of the environment in which this command will be run",
		EnvVar:    "CATALYZE_ENV",
		HideValue: true,
	})
	var settings *models.Settings

	app.Before = func() {
		// TODO auth
		r := config.FileSettingsRetriever{}
		settings = r.GetSettings(*givenEnvName, "", baasHost, paasHost, *username, *password)
		// TODO do this thing
		/*if settings.Pods == nil || len(*settings.Pods) == 0 {
			settings.Pods = helpers.ListPods(settings)
			fmt.Println(settings.Pods)
		}*/
	}
	app.After = func() {
		// TODO save the settings
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
	versionString := fmt.Sprintf("version %s %s\n", config.VERSION, archString)
	app.Version("v version", versionString)
	app.Command("version", "Output the version and quit", func(cmd *cli.Cmd) {
		cmd.Action = app.PrintVersion
	})

	app.Run(os.Args)
}

// InitCLI adds arguments and commands to the given cli instance
func InitCLI(app *cli.Cli, settings *models.Settings) {

	// TODO ideally, we want to upgrade the mow.cli and use the precommand hook to take care
	// of authentication. that way we can create the settings object here and then
	// the commands dont need anyhting but the settings object. then they just have
	// to check if the serviceID or environmentID on the settings object is empty.
	// if required and empty, prompt or throw error as appropriate.

	app.Command(associate.Cmd.Name, associate.Cmd.ShortHelp, associate.Cmd.CmdFunc(settings))
	app.Command(associated.Cmd.Name, associated.Cmd.ShortHelp, associate.Cmd.CmdFunc(settings))
	app.Command(console.Cmd.Name, console.Cmd.ShortHelp, console.Cmd.CmdFunc(settings))
	app.Command(dashboard.Cmd.Name, dashboard.Cmd.ShortHelp, dashboard.Cmd.CmdFunc(settings))
	app.Command(db.Cmd.Name, db.Cmd.ShortHelp, db.Cmd.CmdFunc(settings))
	app.Command(defaultcmd.Cmd.Name, defaultcmd.Cmd.ShortHelp, defaultcmd.Cmd.CmdFunc(settings))
	app.Command(disassociate.Cmd.Name, disassociate.Cmd.ShortHelp, disassociate.Cmd.CmdFunc(settings))
	app.Command(environments.Cmd.Name, environments.Cmd.ShortHelp, environments.Cmd.CmdFunc(settings))
	app.Command(files.Cmd.Name, files.Cmd.ShortHelp, files.Cmd.CmdFunc(settings))
	app.Command(invites.Cmd.Name, invites.Cmd.ShortHelp, invites.Cmd.CmdFunc(settings))
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
