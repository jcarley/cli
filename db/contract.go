package db

import (
	"fmt"
	"os"

	"github.com/catalyzeio/cli/crypto"
	"github.com/catalyzeio/cli/jobs"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/prompts"
	"github.com/catalyzeio/cli/services"
	"github.com/catalyzeio/cli/tasks"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "db",
	ShortHelp: "Tasks for databases",
	LongHelp:  "Tasks for databases",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(BackupSubCmd.Name, BackupSubCmd.ShortHelp, BackupSubCmd.CmdFunc(settings))
			cmd.Command(DownloadSubCmd.Name, DownloadSubCmd.ShortHelp, DownloadSubCmd.CmdFunc(settings))
			cmd.Command(ExportSubCmd.Name, ExportSubCmd.ShortHelp, ExportSubCmd.CmdFunc(settings))
			cmd.Command(ImportSubCmd.Name, ImportSubCmd.ShortHelp, ImportSubCmd.CmdFunc(settings))
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
		}
	},
}

var BackupSubCmd = models.Command{
	Name:      "backup",
	ShortHelp: "Create a new backup",
	LongHelp:  "Create a new backup",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database service to create a backup for (i.e. 'db01')")
			skipPoll := subCmd.BoolOpt("s skip-poll", false, "Whether or not to wait for the backup to finish")
			subCmd.Action = func() {
				err := CmdBackup(*databaseName, *skipPoll, New(settings, crypto.New(), jobs.New(settings)), services.New(settings), tasks.New(settings))
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			subCmd.Spec = "DATABASE_NAME [-s]"
		}
	},
}

var DownloadSubCmd = models.Command{
	Name:      "download",
	ShortHelp: "Download a previously created backup",
	LongHelp:  "Download a previously created backup",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database service which was backed up (i.e. 'db01')")
			backupID := subCmd.StringArg("BACKUP_ID", "", "The ID of the backup to download (found from \"catalyze backup list\")")
			filePath := subCmd.StringArg("FILEPATH", "", "The location to save the downloaded backup to. This location must NOT already exist unless -f is specified")
			force := subCmd.BoolOpt("f force", false, "If a file previously exists at \"filepath\", overwrite it and download the backup")
			subCmd.Action = func() {
				err := CmdDownload(*databaseName, *backupID, *filePath, *force, New(settings, crypto.New(), jobs.New(settings)), prompts.New(), services.New(settings))
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			subCmd.Spec = "DATABASE_NAME BACKUP_ID FILEPATH [-f]"
		}
	},
}

var ExportSubCmd = models.Command{
	Name:      "export",
	ShortHelp: "Export data from a database",
	LongHelp:  "Export data from a database",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database to export data from (i.e. 'db01')")
			filePath := subCmd.StringArg("FILEPATH", "", "The location to save the exported data. This location must NOT already exist unless -f is specified")
			force := subCmd.BoolOpt("f force", false, "If a file previously exists at `filepath`, overwrite it and export data")
			subCmd.Action = func() {
				err := CmdExport(*databaseName, *filePath, *force, New(settings, crypto.New(), jobs.New(settings)), prompts.New(), services.New(settings), tasks.New(settings))
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			subCmd.Spec = "DATABASE_NAME FILEPATH [-f]"
		}
	},
}

var ImportSubCmd = models.Command{
	Name:      "import",
	ShortHelp: "Import data into a database",
	LongHelp:  "Import data into a database",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database to import data to (i.e. 'db01')")
			filePath := subCmd.StringArg("FILEPATH", "", "The location of the file to import to the database")
			mongoCollection := subCmd.StringOpt("c mongo-collection", "", "If importing into a mongo service, the name of the collection to import into")
			mongoDatabase := subCmd.StringOpt("d mongo-database", "", "If importing into a mongo service, the name of the database to import into")
			subCmd.Action = func() {
				err := CmdImport(*databaseName, *filePath, *mongoCollection, *mongoDatabase, New(settings, crypto.New(), jobs.New(settings)), services.New(settings), tasks.New(settings))
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			subCmd.Spec = "DATABASE_NAME FILEPATH [-d [-c]]"
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List created backups",
	LongHelp:  "List created backups",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database service to list backups for (i.e. 'db01')")
			page := subCmd.IntOpt("p page", 1, "The page to view")
			pageSize := subCmd.IntOpt("n page-size", 10, "The number of items to show per page")
			subCmd.Action = func() {
				err := CmdList(*databaseName, *page, *pageSize, New(settings, crypto.New(), jobs.New(settings)), services.New(settings))
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			subCmd.Spec = "DATABASE_NAME [-p] [-n]"
		}
	},
}

// IDb
type IDb interface {
	Backup(service *models.Service) (*models.Task, error)
	Download(backupID, filePath string, service *models.Service) error
	Export(filePath string, task *models.Task, service *models.Service) error
	Import(filePath, mongoCollection, mongoDatabase string, service *models.Service) (*models.Task, error)
	List(page, pageSize int, service *models.Service) (*[]models.Job, error)
	TempUploadURL(service *models.Service) (*models.TempURL, error)
	TempDownloadURL(jobID string, service *models.Service) (*models.TempURL, error)
	TempLogsURL(jobID string, serviceID string) (*models.TempURL, error)
	DumpLogs(taskType string, task *models.Task, service *models.Service) error
}

// SDb is a concrete implementation of IDb
type SDb struct {
	Settings *models.Settings
	Crypto   crypto.ICrypto
	Jobs     jobs.IJobs
}

// New returns an instance of IDb
func New(settings *models.Settings, crypto crypto.ICrypto, jobs jobs.IJobs) IDb {
	return &SDb{
		Settings: settings,
		Crypto:   crypto,
		Jobs:     jobs,
	}
}
