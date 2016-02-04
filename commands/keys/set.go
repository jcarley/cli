package keys

import (
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
	"github.com/mitchellh/go-homedir"
)

var SetSubCmd = models.Command{
	Name:      "set",
	ShortHelp: "Set your auth key",
	LongHelp:  "Set the private key used to sign in instead of username and password. This is expected to correspond to an OpenSSH-formatted RSA public key in the same directory.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			path := cmd.StringArg("PATH_TO_KEY", "", "Relative path to the private key file.")

			cmd.Action = func() {
				fullPath, err := homedir.Expand(*path)
				if err != nil {
					logrus.Fatal(err)
				}

				// make sure both files exist
				_, err = ioutil.ReadFile(fullPath + ".pub")
				if err != nil {
					if os.IsNotExist(err) {
						logrus.Fatalf("Public key file '%s' does not exist.", fullPath+".pub")
					} else {
						logrus.Fatal(err)
					}
				}

				_, err = ioutil.ReadFile(fullPath)
				if err != nil {
					if os.IsNotExist(err) {
						logrus.Fatalf("Private key file '%s' does not exist.", fullPath)
					} else {
						logrus.Fatal(err)
					}
				}

				settings.PrivateKeyPath = fullPath
				settings.SessionToken = ""
				a := auth.New(settings, prompts.New())
				user, err := a.Signin()
				if err != nil {
					logrus.Fatal(err)
				}
				logrus.Infof("Successfully added key and signed in as %s.", user.Email)
				config.SaveSettings(settings)
			}
		}
	},
}
