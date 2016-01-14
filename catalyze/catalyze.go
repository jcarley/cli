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
	"github.com/catalyzeio/cli/models"
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
	/*
		app.Command("files", "Tasks for managing service files", func(cmd *cli.Cmd) {
			cmd.Command("download", "Download a file to your localhost with the same file permissions as on the remote host or print it to stdout", func(subCmd *cli.Cmd) {
				serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to download a file from")
				fileName := subCmd.StringArg("FILE_NAME", "", "The name of the service file from running \"catalyze files list\"")
				output := subCmd.StringOpt("o output", "", "The downloaded file will be saved to the given location with the same file permissions as it has on the remote host. If those file permissions cannot be applied, a warning will be printed and default 0644 permissions applied. If no output is specified, stdout is used.")
				force := subCmd.BoolOpt("f force", false, "If the specified output file already exists, automatically overwrite it")
				subCmd.Action = func() {
					settings := r.GetSettings(true, true, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
					files.DownloadServiceFile(*serviceName, *fileName, *output, *force, settings)
				}
				subCmd.Spec = "SERVICE_NAME FILE_NAME [-o] [-f]"
			})
			cmd.Command("list", "List all files available for a given service", func(subCmd *cli.Cmd) {
				serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to list files for")
				subCmd.Action = func() {
					settings := r.GetSettings(true, true, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
					files.ListServiceFiles(*serviceName, settings)
				}
				subCmd.Spec = "SERVICE_NAME"
			})
		})
		app.Command("invites", "Manage invitations for your environments", func(cmd *cli.Cmd) {
			cmd.Command("list", "List all pending environment invitations", func(subCmd *cli.Cmd) {
				subCmd.Action = func() {
					settings := r.GetSettings(false, false, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
					invites.ListInvites(settings)
				}
			})
			cmd.Command("rm", "Remove a pending environment invitation", func(subCmd *cli.Cmd) {
				inviteID := subCmd.StringArg("INVITE_ID", "", "The ID of an invitation to remove")
				subCmd.Action = func() {
					settings := r.GetSettings(false, false, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
					invites.RmInvite(*inviteID, settings)
				}
				subCmd.Spec = "INVITE_ID"
			})
			cmd.Command("send", "Send an invite to a user by email for the associated environment", func(subCmd *cli.Cmd) {
				email := subCmd.StringArg("EMAIL", "", "The email of a user to invite to the associated environment. This user does not need to have a Catalyze account prior to sending the invitation")
				subCmd.Action = func() {
					settings := r.GetSettings(false, false, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
					invites.InviteUser(*email, settings)
				}
				subCmd.Spec = "EMAIL"
			})
		})
		app.Command("logs", "Show the logs in your terminal streamed from your logging dashboard", func(cmd *cli.Cmd) {
			query := cmd.StringArg("QUERY", "*", "The query to send to your logging dashboard's elastic search (regex is supported)")
			follow := cmd.BoolOpt("f follow", false, "Tail/follow the logs (Equivalent to -t)")
			tail := cmd.BoolOpt("t tail", false, "Tail/follow the logs (Equivalent to -f)")
			hours := cmd.IntOpt("hours", 0, "The number of hours before now (in combination with minutes and seconds) to retrieve logs")
			mins := cmd.IntOpt("minutes", 1, "The number of minutes before now (in combination with hours and seconds) to retrieve logs")
			secs := cmd.IntOpt("seconds", 0, "The number of seconds before now (in combination with hours and minutes) to retrieve logs")
			cmd.Action = func() {
				settings := r.GetSettings(false, false, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
				logs.Logs(*query, *tail || *follow, *hours, *mins, *secs, settings)
			}
			cmd.Spec = "[QUERY] [(-f | -t)] [--hours] [--minutes] [--seconds]"
		})
		app.Command("logout", "Clear the stored user information from your local machine", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				settings := r.GetSettings(false, false, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
				logout.Logout(settings)
			}
		})
		app.Command("metrics", "Print service and environment metrics in your local time zone", func(cmd *cli.Cmd) {
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to print metrics for")
			json := cmd.BoolOpt("json", false, "Output the data as json")
			csv := cmd.BoolOpt("csv", false, "Output the data as csv")
			spark := cmd.BoolOpt("spark", false, "Output the data using spark lines")
			stream := cmd.BoolOpt("stream", false, "Repeat calls once per minute until this process is interrupted.")
			mins := cmd.IntOpt("m mins", 1, "How many minutes worth of logs to retrieve.")
			cmd.Action = func() {
				settings := r.GetSettings(true, true, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
				metrics.Metrics(*serviceName, *json, *csv, *spark, *stream, *mins, settings)
			}
			cmd.Spec = "[SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]"
		})
		app.Command("rake", "Execute a rake task", func(cmd *cli.Cmd) {
			taskName := cmd.StringArg("TASK_NAME", "", "The name of the rake task to run")
			cmd.Action = func() {
				settings := r.GetSettings(true, true, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
				rake.Rake(*taskName, settings)
			}
			cmd.Spec = "TASK_NAME"
		})
		app.Command("redeploy", "Redeploy a service without having to do a git push", func(cmd *cli.Cmd) {
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to redeploy (i.e. 'app01')")
			cmd.Action = func() {
				settings := r.GetSettings(true, true, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
				redeploy.Redeploy(*serviceName, settings)
			}
			cmd.Spec = "SERVICE_NAME"
		})
		app.Command("services", "List all services for your environment", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				settings := r.GetSettings(true, true, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
				services.ListServices(settings)
			}
		})
		app.Command("ssl", "Perform operations on local certificates to verify their validity", func(cmd *cli.Cmd) {
			cmd.Command("verify", "Verify whether a certificate chain is complete and if it matches the given private key", func(subCmd *cli.Cmd) {
				chain := subCmd.StringArg("CHAIN", "", "The path to your full certificate chain in PEM format")
				privateKey := subCmd.StringArg("PRIVATE_KEY", "", "The path to your private key in PEM format")
				hostname := subCmd.StringArg("HOSTNAME", "", "The hostname that should match your certificate (i.e. \"*.catalyze.io\")")
				selfSigned := subCmd.BoolOpt("s self-signed", false, "Whether or not the certificate is self signed. If set, chain verification is skipped")
				subCmd.Action = func() {
					ssl.VerifyChain(*chain, *privateKey, *hostname, *selfSigned)
				}
				subCmd.Spec = "CHAIN PRIVATE_KEY HOSTNAME [-s]"
			})
		})
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
