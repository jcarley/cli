package keys

import (
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	libKeys "github.com/catalyzeio/cli/lib/keys"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
	"github.com/mitchellh/go-homedir"
)

var AddSubCmd = models.Command{
	Name:      "add",
	ShortHelp: "Add a public key",
	LongHelp:  "Add a new RSA public key to your own user account",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			name := cmd.StringArg("NAME", "", "The name for the new key, for your own purposes.")
			path := cmd.StringArg("PATH_TO_KEY", "", "Relative path to the public key file.")

			cmd.Action = func() {
				fullPath, err := homedir.Expand(*path)
				if err != nil {
					logrus.Fatal(err)
				}

				keyBytes, err := ioutil.ReadFile(fullPath)
				if err != nil {
					logrus.Fatal(err)
				}
				err = libKeys.Add(settings, *name, string(keyBytes))
				if err != nil {
					logrus.Fatal(err)
				}

				logrus.Printf("Key '%s' added to your account.", *name)
			}
		}
	},
}
