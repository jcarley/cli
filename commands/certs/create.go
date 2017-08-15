package certs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/commands/ssl"
	"github.com/daticahealth/cli/config"
	"github.com/daticahealth/cli/models"
)

func CmdCreate(name, pubKeyPath, privKeyPath string, selfSigned, resolve, letsEncrypt bool, ic ICerts, is services.IServices, issl ssl.ISSL) error {
	if strings.ContainsAny(name, config.InvalidChars) {
		return fmt.Errorf("Invalid cert name. Names must not contain the following characters: %s", config.InvalidChars)
	}
	service, err := is.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}
	if letsEncrypt {
		err = ic.CreateLetsEncrypt(name, service.ID)
		if err != nil {
			return err
		}
	} else {
		if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
			return fmt.Errorf("A cert does not exist at path '%s'", pubKeyPath)
		}
		if _, err := os.Stat(privKeyPath); os.IsNotExist(err) {
			return fmt.Errorf("A private key does not exist at path '%s'", privKeyPath)
		}
		err := issl.Verify(pubKeyPath, privKeyPath, name, selfSigned)
		var pubKeyBytes []byte
		var privKeyBytes []byte
		if err != nil && !ssl.IsHostnameMismatchErr(err) {
			if ssl.IsIncompleteChainErr(err) && resolve {
				pubKeyBytes, err = issl.Resolve(pubKeyPath)
				if err != nil {
					return fmt.Errorf("Could not resolve the incomplete certificate chain. If this is a self signed certificate, please re-run this command with the '-s' option: %s", err.Error())
				}
			} else {
				return err
			}
		}
		if pubKeyBytes == nil {
			pubKeyBytes, err = ioutil.ReadFile(pubKeyPath)
			if err != nil {
				return err
			}
		}
		if privKeyBytes == nil {
			privKeyBytes, err = ioutil.ReadFile(privKeyPath)
			if err != nil {
				return err
			}
		}
		err = ic.Create(name, string(pubKeyBytes), string(privKeyBytes), service.ID)
		if err != nil {
			return err
		}
	}
	logrus.Printf("Created '%s'", name)
	logrus.Println("To make use of your cert, you need to add a site with the \"datica sites create\" command")
	return nil
}

func (c *SCerts) Create(name, pubKey, privKey, svcID string) error {
	cert := models.Cert{
		Name:    name,
		PubKey:  pubKey,
		PrivKey: privKey,
	}
	b, err := json.Marshal(cert)
	if err != nil {
		return err
	}
	headers := c.Settings.HTTPManager.GetHeaders(c.Settings.SessionToken, c.Settings.Version, c.Settings.Pod, c.Settings.UsersID)
	resp, statusCode, err := c.Settings.HTTPManager.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/certs", c.Settings.PaasHost, c.Settings.PaasHostVersion, c.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return err
	}
	return c.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}

func (c *SCerts) CreateLetsEncrypt(name, svcID string) error {
	var cert = struct {
		Name        string `json:"name"`
		LetsEncrypt bool   `json:"letsEncrypt"`
	}{
		Name:        name,
		LetsEncrypt: true,
	}
	b, err := json.Marshal(cert)
	if err != nil {
		return err
	}
	headers := c.Settings.HTTPManager.GetHeaders(c.Settings.SessionToken, c.Settings.Version, c.Settings.Pod, c.Settings.UsersID)
	resp, statusCode, err := c.Settings.HTTPManager.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/certs", c.Settings.PaasHost, c.Settings.PaasHostVersion, c.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return err
	}
	return c.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
