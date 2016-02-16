package keys

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/httpclient"
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
				err := CmdRemove(New(settings), *name)
				if err != nil {
					logrus.Fatal(err)
				}
				logrus.Printf("Key '%s' has been removed from your account.", *name)
			}
		}
	},
}

func CmdRemove(k IKeys, name string) error {
	return k.Remove(name)
}

func (k *SKeys) Remove(name string) error {
	headers := httpclient.GetHeaders(k.Settings.SessionToken, k.Settings.Version, k.Settings.Pod)
	resp, status, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/keys/%s", k.Settings.AuthHost, k.Settings.AuthHostVersion, name), headers)
	if err != nil {
		return err
	}
	if httpclient.IsError(status) {
		return httpclient.ConvertError(resp, status)
	}
	return nil
}
