package keys

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/httpclient"
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
				err := CmdList(New(settings), *printKeys)
				if err != nil {
					logrus.Fatal(err)
				}
			}
		}
	},
}

func CmdList(k IKeys, printKeys bool) error {
	keys, err := k.List()
	if err != nil {
		return err
	}

	for _, key := range keys {
		if printKeys {
			logrus.Printf("%s = %s\n", key.Name, key.Key)
		} else {
			logrus.Printf("%s", key.Name)
		}
	}
	return nil
}

func (k *SKeys) List() ([]models.UserKey, error) {
	headers := httpclient.GetHeaders(k.Settings.SessionToken, k.Settings.Version, k.Settings.Pod)
	resp, status, err := httpclient.Get(nil, fmt.Sprintf("%s%s/keys", k.Settings.AuthHost, k.Settings.AuthHostVersion), headers)
	if err != nil {
		return nil, err
	}

	keys := []models.UserKey{}
	err = httpclient.ConvertResp(resp, status, &keys)
	return keys, err
}
