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
	"github.com/catalyzeio/cli/updater"
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
	/*
		app.Command("status", "Get quick readout of the current status of your associated environment and all of its services", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				settings := r.GetSettings(true, true, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
				status.Status(settings)
			}
		})
		app.Command("support-ids", "Print out various IDs related to your associated environment to be used when contacting Catalyze support", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				settings := r.GetSettings(true, true, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
				supportids.SupportIds(settings)
			}
		})
		app.Command("update", "Checks for available updates and updates the CLI if a new update is available", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				update.Update()
			}
		})
		app.Command("vars", "Interaction with environment variables for the associated environment", func(cmd *cli.Cmd) {
			cmd.Command("list", "List all environment variables", func(subCmd *cli.Cmd) {
				subCmd.Action = func() {
					settings := r.GetSettings(true, true, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
					vars.ListVars(settings)
				}
			})
			cmd.Command("set", "Set one or more new environment variables or update the values of existing ones", func(subCmd *cli.Cmd) {
				variables := subCmd.Strings(cli.StringsOpt{
					Name:      "v variable",
					Value:     []string{},
					Desc:      "The env variable to set or update in the form \"<key>=<value>\"",
					HideValue: true,
				})
				subCmd.Action = func() {
					settings := r.GetSettings(true, true, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
					vars.SetVar(*variables, settings)
				}
				subCmd.Spec = "-v..."
			})
			cmd.Command("unset", "Unset (delete) an existing environment variable", func(subCmd *cli.Cmd) {
				variable := subCmd.StringArg("VARIABLE", "", "The name of the environment variable to unset")
				subCmd.Action = func() {
					settings := r.GetSettings(true, true, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
					vars.UnsetVar(*variable, settings)
				}
				subCmd.Spec = "VARIABLE"
			})
		})
		app.Command("whoami", "Retrieve your user ID", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				settings := r.GetSettings(false, false, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
				whoami.WhoAmI(settings)
			}
		})
		app.Command("worker", "Start a background worker", func(cmd *cli.Cmd) {
			target := cmd.StringArg("TARGET", "", "The name of the Procfile target to invoke as a worker")
			cmd.Action = func() {
				settings := r.GetSettings(true, true, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
				worker.Worker(*target, settings)
			}
			cmd.Spec = "TARGET"
		})*/
}
