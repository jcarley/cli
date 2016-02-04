package keys

import (
	"github.com/Sirupsen/logrus"
	libKeys "github.com/catalyzeio/cli/lib/keys"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

var RemoveSubCmd = models.Command{
	Name:      "rm",
	ShortHelp: "Remove a public key",
	LongHelp:  "Remove a public key from your own account, by name.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			name := cmd.StringArg("NAME", "", "The name of the key to remove.")

			cmd.Action = func() {
				err := libKeys.Remove(settings, *name)
				if err != nil {
					logrus.Fatal(err)
				}

				logrus.Printf("Key '%s' has been removed from your account.", *name)
			}
		}
	},
}
