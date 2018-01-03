package certs

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/models"
	"github.com/olekukonko/tablewriter"
)

func CmdList(ic ICerts, is services.IServices, downStream string) error {
	service, err := is.RetrieveByLabel(downStream)
	if err != nil {
		return err
	}
	certs, err := ic.List(service.ID)
	if err != nil {
		return err
	}
	if certs == nil || len(*certs) == 0 {
		logrus.Println("No certs found")
		return nil
	}

	data := [][]string{{"NAME", "LET'S ENCRYPT STATUS"}}
	for _, cert := range *certs {
		data = append(data, []string{cert.Name, cert.LetsEncrypt.String()})
	}

	table := tablewriter.NewWriter(logrus.StandardLogger().Out)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetAutoWrapText(false)
	table.AppendBulk(data)
	table.Render()
	return nil
}

func (c *SCerts) List(svcID string) (*[]models.Cert, error) {
	headers := c.Settings.HTTPManager.GetHeaders(c.Settings.SessionToken, c.Settings.Version, c.Settings.Pod, c.Settings.UsersID)
	resp, statusCode, err := c.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/certs", c.Settings.PaasHost, c.Settings.PaasHostVersion, c.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return nil, err
	}
	var certs []models.Cert
	err = c.Settings.HTTPManager.ConvertResp(resp, statusCode, &certs)
	if err != nil {
		return nil, err
	}
	return &certs, nil
}
