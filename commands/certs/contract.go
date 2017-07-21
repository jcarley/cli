package certs

import (
	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/commands/ssl"
	"github.com/daticahealth/cli/config"
	"github.com/daticahealth/cli/lib/auth"
	"github.com/daticahealth/cli/lib/prompts"
	"github.com/daticahealth/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "certs",
	ShortHelp: "Manage your SSL certificates and domains",
	LongHelp:  "The `certs` command gives access to certificate and private key management for public facing services. The certs command cannot be run directly but has sub commands.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.CommandLong(CreateSubCmd.Name, CreateSubCmd.ShortHelp, CreateSubCmd.LongHelp, CreateSubCmd.CmdFunc(settings))
			cmd.CommandLong(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.LongHelp, ListSubCmd.CmdFunc(settings))
			cmd.CommandLong(RmSubCmd.Name, RmSubCmd.ShortHelp, RmSubCmd.LongHelp, RmSubCmd.CmdFunc(settings))
			cmd.CommandLong(UpdateSubCmd.Name, UpdateSubCmd.ShortHelp, UpdateSubCmd.LongHelp, UpdateSubCmd.CmdFunc(settings))
		}
	},
}

var CreateSubCmd = models.Command{
	Name:      "create",
	ShortHelp: "Create a new domain with an SSL certificate and private key or create a Let's Encrypt certificate",
	LongHelp: "`certs create` allows you to upload an SSL certificate and private key which can be used to secure your public facing code service. " +
		"Alternatively, you may opt to create a Let's Encrypt certificate. When creating a Let's Encrypt certificate, you only need to provide the certificate name along with the \"-l\" flag. " +
		"Let's Encrypt certificates are issued asynchronously and may not be available immediately. Use the [certs list](#certs-list) command to check on the issuance status. " +
		"Once issued, Let's Encrypt certificates automatically renew before expiring. " +
		"Cert creation can be done at any time, even after environment provisioning, but must be done before [creating a site](#sites-create). " +
		"When uploading a custom cert, the CLI will check to ensure the certificate and private key match. If you are using a self signed cert, pass in the `-s` flag and the hostname check will be skipped. " +
		"Datica requires that your certificate include your own certificate, intermediate certificates, and the root certificate in that order. " +
		"If you only include your certificate, the CLI will attempt to resolve this and fetch intermediate and root certificates for you. " +
		"It is advised that you create a full chain before running this command as the `-r` flag is accomplished on a \"best effort\" basis.\n\n" +
		"Here are a few sample commands\n\n" +
		"```\ndatica -E \"<your_env_name>\" certs create wildcard_mysitecom ~/path/to/cert.pem ~/path/to/priv.key\n" +
		"datica -E \"<your_env_name>\" certs create my.site.com --lets-encrypt\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			name := subCmd.StringArg("NAME", "", "The name of this SSL certificate plus private key pair")
			pubKeyPath := subCmd.StringArg("PUBLIC_KEY_PATH", "", "The path to a public key file in PEM format")
			privKeyPath := subCmd.StringArg("PRIVATE_KEY_PATH", "", "The path to an unencrypted private key file in PEM format")
			selfSigned := subCmd.BoolOpt("s self-signed", false, "Whether or not the given SSL certificate and private key are self signed")
			resolve := subCmd.BoolOpt("r resolve", true, "Whether or not to attempt to automatically resolve incomplete SSL certificate issues")
			letsEncrypt := subCmd.BoolOpt("l lets-encrypt", false, "Whether or not this is a Let's Encrypt certificate")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdCreate(*name, *pubKeyPath, *privKeyPath, *selfSigned, *resolve, *letsEncrypt, New(settings), services.New(settings), ssl.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "NAME ((PUBLIC_KEY_PATH PRIVATE_KEY_PATH [-s] [-r]) | -l)"
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all existing domains that have SSL certificate and private key pairs",
	LongHelp: "`certs list` lists all of the available certs you have created on your environment. " +
		"The displayed names are the names that should be used as the `CERT_NAME` parameter in the [sites create](#sites-create) command. " +
		"If any certs are Let's Encrypt certs, the issuance status will also be shown. " +
		"Here is a sample command\n\n" +
		"```\ndatica -E \"<your_env_name>\" certs list\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(settings); err != nil {
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
	ShortHelp: "Remove an existing domain and its associated SSL certificate and private key pair",
	LongHelp: "`certs rm` allows you to delete old certificate and private key pairs. Only certs that are not in use by a site can be deleted. Here is a sample command\n\n" +
		"```\ndatica -E \"<your_env_name>\" certs rm mywebsite.com\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			name := subCmd.StringArg("NAME", "", "The name of the certificate to remove")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(settings); err != nil {
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

var UpdateSubCmd = models.Command{
	Name:      "update",
	ShortHelp: "Update the SSL certificate and private key pair for an existing domain",
	LongHelp: "`certs update` works nearly identical to the [certs create](#certs-create) command. " +
		"All rules regarding self signed certs and certificate resolution from the `certs create` command apply to the `certs update` command. " +
		"Let's Encrypt certs cannot be updated since they are automatically renewed before expiring. " +
		"This is useful for when your certificates have expired and you need to upload new ones. Update your certs and then redeploy your service_proxy. Here is a sample command\n\n" +
		"```\ndatica -E \"<your_env_name>\" certs update mywebsite.com ~/path/to/new/cert.pem ~/path/to/new/priv.key\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			name := subCmd.StringArg("NAME", "", "The name of this SSL certificate and private key pair")
			pubKeyPath := subCmd.StringArg("PUBLIC_KEY_PATH", "", "The path to a public key file in PEM format")
			privKeyPath := subCmd.StringArg("PRIVATE_KEY_PATH", "", "The path to an unencrypted private key file in PEM format")
			selfSigned := subCmd.BoolOpt("s self-signed", false, "Whether or not the given SSL certificate and private key are self signed")
			resolve := subCmd.BoolOpt("r resolve", true, "Whether or not to attempt to automatically resolve incomplete SSL certificate issues")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdUpdate(*name, *pubKeyPath, *privKeyPath, *selfSigned, *resolve, New(settings), services.New(settings), ssl.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "NAME PUBLIC_KEY_PATH PRIVATE_KEY_PATH [-s] [-r]"
		}
	},
}

// ICerts
type ICerts interface {
	Create(name, pubKey, privKey, svcID string) error
	CreateLetsEncrypt(name, svcID string) error
	Update(name, pubKey, privKey, svcID string) error
	List(svcID string) (*[]models.Cert, error)
	Rm(name, svcID string) error
}

// SCerts is a concrete implementation of ICerts
type SCerts struct {
	Settings *models.Settings
}

// New returns an instance of ICerts
func New(settings *models.Settings) ICerts {
	return &SCerts{
		Settings: settings,
	}
}
