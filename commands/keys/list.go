package keys

import (
	"github.com/Sirupsen/logrus"
	libKeys "github.com/catalyzeio/cli/lib/keys"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List your public keys",
	LongHelp:  "List the names of all public keys currently attached to your user",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			printKeys := cmd.BoolOpt("include-keys", false, "Print out the values of the public keys, as well as names.")

			cmd.Action = func() {
				keys, err := libKeys.List(settings)
				if err != nil {
					logrus.Fatal(err.Error())
				}

				for _, key := range keys {
					if *printKeys {
						logrus.Printf("%s = %s\n", key.Name, key.Key)
					} else {
						logrus.Printf("%s", key.Name)
					}
				}
			}
		}
	},
}
