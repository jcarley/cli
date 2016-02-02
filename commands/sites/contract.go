package sites

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/files"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/commands/ssl"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "sites",
	ShortHelp: "Tasks for updating sites, including hostnames, SSL certificates, and private keys",
	LongHelp:  "Tasks for updating sites, including hostnames, SSL certificates, and private keys",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(CreateSubCmd.Name, CreateSubCmd.ShortHelp, CreateSubCmd.CmdFunc(settings))
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(RmSubCmd.Name, RmSubCmd.ShortHelp, RmSubCmd.CmdFunc(settings))
		}
	},
}

var CreateSubCmd = models.Command{
	Name:      "create",
	ShortHelp: "Create a new site including a hostname, SSL certificate, and private key",
	LongHelp:  "Create a new site including a hostname, SSL certificate, and private key",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			chain := subCmd.StringArg("CHAIN", "", "The path to your full certificate chain in PEM format")
			privateKey := subCmd.StringArg("PRIVATE_KEY", "", "The path to your private key in PEM format")
			hostname := subCmd.StringArg("HOSTNAME", "", "The hostname of your service")
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to add this site configuration to (i.e. 'app01')")
			wildcard := subCmd.BoolOpt("w wildcard", false, "Whether or not the given SSL certificate is for a wildcard domain")
			selfSigned := subCmd.BoolOpt("s self-signed", false, "Whether or not the given SSL certificate is a self-signed cert")
			subCmd.Action = func() {
				err := CmdCreate(*hostname, *chain, *privateKey, *serviceName, *wildcard, *selfSigned, New(settings), ssl.New(settings), files.New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "CHAIN PRIVATE_KEY HOSTNAME SERVICE_NAME [-w] [-s]"
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all site configurations for all code services",
	LongHelp:  "List all site configurations for all code services",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				err := CmdList(New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

var RmSubCmd = models.Command{
	Name:      "rm",
	ShortHelp: "Remove a site configuration",
	LongHelp:  "Remove a site configuration",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			siteID := subCmd.IntArg("SITE_ID", 0, "The id of the site configuration to delete")
			subCmd.Action = func() {
				err := CmdRm(*siteID, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "SITE_ID"
		}
	},
}

// ISites
type ISites interface {
	Create(svcID string, site *models.Site) error
	List(svcID string) (*[]models.Site, error)
	Rm(svcID string, siteID int) error
}

// SSites is a concrete implementation of ISites
type SSites struct {
	Settings *models.Settings
}

// New returns an instance of ISites
func New(settings *models.Settings) ISites {
	return &SSites{
		Settings: settings,
	}
}
