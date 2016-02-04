package keys

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/httpclient"
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
				err := CmdAdd(New(settings), *name, *path)
				if err != nil {
					logrus.Fatal(err)
				}
				logrus.Printf("Key '%s' added to your account.", *name)
			}
		}
	},
}

func CmdAdd(k IKeys, name string, path string) error {
	fullPath, err := homedir.Expand(path)
	if err != nil {
		return err
	}
	keyBytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	err = k.Add(name, string(keyBytes))
	if err != nil {
		return err
	}
	return nil
}

// Add adds a new public key to the authenticated user's account
func (k *SKeys) Add(name string, publicKey string) error {
	body, err := json.Marshal(struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	}{
		Key:  publicKey,
		Name: name,
	})
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(k.Settings.SessionToken, k.Settings.Version, k.Settings.Pod)
	resp, status, err := httpclient.Post(body, fmt.Sprintf("%s%s/keys", k.Settings.AuthHost, k.Settings.AuthHostVersion), headers)
	if err != nil {
		return err
	}
	if httpclient.IsError(status) {
		return httpclient.ConvertError(resp, status)
	}
	return nil
}
