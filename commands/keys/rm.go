package keys

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/deploykeys"
	"github.com/daticahealth/cli/config"
)

func CmdRemove(name, privateKeyPath string, ik IKeys, id deploykeys.IDeployKeys) error {
	if privateKeyPath != "" {
		// if ssh key auth is being used, don't let the key used for key auth be removed
		b, err := ioutil.ReadFile(privateKeyPath)
		if err != nil {
			return err
		}
		privKey, err := id.ParsePrivateKey(b)
		if err != nil {
			return err
		}
		pubKey, err := id.ExtractPublicKey(privKey)
		if err != nil {
			return err
		}
		pubKeyString := string(ssh.MarshalAuthorizedKey(pubKey))
		userKeys, err := ik.List()
		if err != nil {
			return err
		}
		for _, uk := range *userKeys {
			if strings.TrimSpace(pubKeyString) == strings.TrimSpace(uk.Key) {
				return errors.New("You cannot remove the key that is currently being used for authentication. Run \"datica clear --private-key\" to remove your SSH key authentication settings.")
			}
		}
	}
	if strings.ContainsAny(name, config.InvalidChars) {
		return fmt.Errorf("Invalid key name. Names must not contain the following characters: %s", config.InvalidChars)
	}
	err := ik.Remove(name)
	if err != nil {
		return err
	}
	logrus.Printf("Key '%s' has been removed from your account.", name)
	return nil
}

func (k *SKeys) Remove(name string) error {
	headers := k.Settings.HTTPManager.GetHeaders(k.Settings.SessionToken, k.Settings.Version, k.Settings.Pod, k.Settings.UsersID)
	resp, status, err := k.Settings.HTTPManager.Delete(nil, fmt.Sprintf("%s%s/keys/%s", k.Settings.AuthHost, k.Settings.AuthHostVersion, name), headers)
	if err != nil {
		return err
	}
	return k.Settings.HTTPManager.ConvertResp(resp, status, nil)
}
