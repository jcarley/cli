package services

import (
	"fmt"
	"os"

	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "services",
	ShortHelp: "List all services for your environment",
	LongHelp:  "List all services for your environment",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				is := New(settings, "", "")
				err := CmdServices(is)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
		}
	},
}

// IServices
type IServices interface {
	List() (*[]models.Service, error)
	Retrieve() (*models.Service, error)
	RetrieveByLabel() (*models.Service, error)
}

// SServices is a concrete implementation of IServices
type SServices struct {
	Settings *models.Settings

	SvcID    string
	SvcLabel string
}

// New generates a new instance of IServices
func New(settings *models.Settings, svcID string, svcLabel string) IServices {
	return &SServices{
		Settings: settings,

		SvcID:    svcID,
		SvcLabel: svcLabel,
	}
}
