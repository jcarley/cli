package certs

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/config"
)

func CmdRm(name string, ic ICerts, is services.IServices, downStream string) error {
	if strings.ContainsAny(name, config.InvalidChars) {
		return fmt.Errorf("Invalid cert name. Names must not contain the following characters: %s", config.InvalidChars)
	}
	service, err := is.RetrieveByLabel(downStream)
	if err != nil {
		return err
	}
	err = ic.Rm(name, service.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Removed '%s'", name)
	return nil
}

func (c *SCerts) Rm(name, svcID string) error {
	headers := c.Settings.HTTPManager.GetHeaders(c.Settings.SessionToken, c.Settings.Version, c.Settings.Pod, c.Settings.UsersID)
	resp, statusCode, err := c.Settings.HTTPManager.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/certs/%s", c.Settings.PaasHost, c.Settings.PaasHostVersion, c.Settings.EnvironmentID, svcID, name), headers)
	if err != nil {
		return err
	}
	return c.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
