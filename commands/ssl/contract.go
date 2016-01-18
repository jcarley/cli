package ssl

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "ssl",
	ShortHelp: "Perform operations on local certificates to verify their validity",
	LongHelp:  "Perform operations on local certificates to verify their validity",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(VerifySubCmd.Name, VerifySubCmd.ShortHelp, VerifySubCmd.CmdFunc(settings))
		}
	},
}

var VerifySubCmd = models.Command{
	Name:      "verify",
	ShortHelp: "Verify whether a certificate chain is complete and if it matches the given private key",
	LongHelp:  "Verify whether a certificate chain is complete and if it matches the given private key",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			chain := subCmd.StringArg("CHAIN", "", "The path to your full certificate chain in PEM format")
			privateKey := subCmd.StringArg("PRIVATE_KEY", "", "The path to your private key in PEM format")
			hostname := subCmd.StringArg("HOSTNAME", "", "The hostname that should match your certificate (i.e. \"*.catalyze.io\")")
			selfSigned := subCmd.BoolOpt("s self-signed", false, "Whether or not the certificate is self signed. If set, chain verification is skipped")
			subCmd.Action = func() {
				err := CmdVerify(*chain, *privateKey, *hostname, *selfSigned, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "CHAIN PRIVATE_KEY HOSTNAME [-s]"
		}
	},
}

// ISSL
type ISSL interface {
	Verify(chainPath, privateKeyPath, hostname string, selfSigned bool) error
}

// SSSL is a concrete implementation of ISSL
type SSSL struct {
	Settings *models.Settings
}

// New generates a new instance of ISSL
func New(settings *models.Settings) ISSL {
	return &SSSL{
		Settings: settings,
	}
}
