package certs

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
)

func CmdRm(hostname string, ic ICerts, is services.IServices) error {
	if strings.ContainsAny(hostname, config.InvalidChars) {
		return fmt.Errorf("Invalid cert hostname. Hostnames must not contain the following characters: %s", config.InvalidChars)
	}
	service, err := is.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}
	err = ic.Rm(hostname, service.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Removed '%s'", hostname)
	return nil
}

func (c *SCerts) Rm(hostname, svcID string) error {
	headers := c.Settings.HTTPManager.GetHeaders(c.Settings.SessionToken, c.Settings.Version, c.Settings.Pod, c.Settings.UsersID)
	resp, statusCode, err := c.Settings.HTTPManager.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/certs/%s", c.Settings.PaasHost, c.Settings.PaasHostVersion, c.Settings.EnvironmentID, svcID, hostname), headers)
	if err != nil {
		return err
	}
	return c.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
