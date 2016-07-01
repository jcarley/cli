package sites

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
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
			cmd.Command(ShowSubCmd.Name, ShowSubCmd.ShortHelp, ShowSubCmd.CmdFunc(settings))
		}
	},
}

var CreateSubCmd = models.Command{
	Name:      "create",
	ShortHelp: "Create a new site linking it to an existing cert instance",
	LongHelp:  "Create a new site linking it to an existing cert instance",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			name := subCmd.StringArg("SITE_NAME", "", "The name of the site to be created. This will be used in this site's nginx configuration file (i.e. \".example.com\")")
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to add this site configuration to (i.e. 'app01')")
			hostname := subCmd.StringArg("HOSTNAME", "", "The hostname used in the creation of a certs instance with the 'certs' command (i.e. \"star_example_com\")")
			clientMaxBodySize := subCmd.IntOpt("client-max-body-size", -1, "The 'client_max_body_size' nginx config specified in megabytes")
			proxyConnectTimeout := subCmd.IntOpt("proxy-connect-timeout", -1, "The 'proxy_connect_timeout' nginx config specified in seconds")
			proxyReadTimeout := subCmd.IntOpt("proxy-read-timeout", -1, "The 'proxy_read_timeout' nginx config specified in seconds")
			proxySendTimeout := subCmd.IntOpt("proxy-send-timeout", -1, "The 'proxy_send_timeout' nginx config specified in seconds")
			proxyUpstreamTimeout := subCmd.IntOpt("proxy-upstream-timeout", -1, "The 'proxy_next_upstream_timeout' nginx config specified in seconds")
			enableCORS := subCmd.BoolOpt("enable-cors", false, "Enable or disable all features related to full CORS support")
			enableWebSockets := subCmd.BoolOpt("enable-websockets", false, "Enable or disable all features related to full websockets support")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdCreate(*name, *serviceName, *hostname, *clientMaxBodySize, *proxyConnectTimeout, *proxyReadTimeout, *proxySendTimeout, *proxyUpstreamTimeout, *enableCORS, *enableWebSockets, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "NAME SERVICE_NAME HOSTNAME [--client-max-body-size] [--proxy-connect-timeout] [--proxy-read-timeout] [--proxy-send-timeout] [--proxy-upstream-timeout] [--enable-cors] [--enable-web-sockets]"
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List details for all site configurations",
	LongHelp:  "List details for all site configurations",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
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
			name := subCmd.StringArg("NAME", "", "The name of the site configuration to delete")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRm(*name, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "NAME"
		}
	},
}

var ShowSubCmd = models.Command{
	Name:      "show",
	ShortHelp: "Shows the details for a given site",
	LongHelp:  "Shows the details for a given site",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			name := subCmd.StringArg("NAME", "", "The name of the site configuration to show")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdShow(*name, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "NAME"
		}
	},
}

// ISites
type ISites interface {
	Create(name, cert, upstreamServiceID, svcID string, siteValues map[string]interface{}) (*models.Site, error)
	List(svcID string) (*[]models.Site, error)
	Retrieve(siteID int, svcID string) (*models.Site, error)
	Rm(siteID int, svcID string) error
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
