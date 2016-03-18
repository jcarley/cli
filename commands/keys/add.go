package keys

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/deploykeys"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/mitchellh/go-homedir"
)

func CmdAdd(name, path string, ik IKeys, id deploykeys.IDeployKeys) error {
	if strings.ContainsAny(name, config.InvalidChars) {
		return fmt.Errorf("Invalid key name. Names must not contain the following characters: %s", config.InvalidChars)
	}
	fullPath, err := homedir.Expand(path)
	if err != nil {
		return err
	}
	keyBytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	k, err := id.ParsePublicKey(keyBytes)
	if err != nil {
		return err
	}
	key := ssh.MarshalAuthorizedKey(k)
	if err != nil {
		return err
	}
	err = ik.Add(name, string(key))
	if err != nil {
		return err
	}
	logrus.Printf("Key '%s' added to your account.", name)
	return nil
}

// Add adds a new public key to the authenticated user's account
func (k *SKeys) Add(name, publicKey string) error {
	body, err := json.Marshal(models.UserKey{
		Key:  publicKey,
		Name: name,
	})
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(k.Settings.SessionToken, k.Settings.Version, k.Settings.Pod, k.Settings.UsersID)
	resp, status, err := httpclient.Post(body, fmt.Sprintf("%s%s/keys", k.Settings.AuthHost, k.Settings.AuthHostVersion), headers)
	if err != nil {
		return err
	}
	if httpclient.IsError(status) {
		return httpclient.ConvertError(resp, status)
	}
	return nil
}
