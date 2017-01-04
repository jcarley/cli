package maintenance

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/olekukonko/tablewriter"
)

func CmdShow(svcName, envID, podID string, im IMaintenance, is services.IServices) error {
	svcs, err := is.ListByEnvID(envID, podID)
	if err != nil {
		return err
	}
	if svcs == nil || len(*svcs) == 0 {
		logrus.Println("No services found")
		return nil
	}
	serviceProxy, err := is.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}
	svcMaintenance, err := im.List(serviceProxy.ID)
	if err != nil {
		return err
	}

	data := [][]string{{"SERVICE", "MAINTENANCE MODE", "ENABLED AT"}}
	for _, svc := range *svcs {
		if svc.Type == "code" && (svcName == "" || svc.Label == svcName) {
			createdAt := ""
			status := "disabled"
			for _, mm := range *svcMaintenance {
				if mm.UpstreamID == svc.ID {
					createdAt = mm.CreatedAt
					status = "enabled"
				}
			}
			data = append(data, []string{svc.Label, status, createdAt})
		}
	}
	if len(data) == 1 {
		if svcName == "" {
			logrus.Println("No code services found")
		} else {
			logrus.Printf("No code service found with the label %s", svcName)
		}
		return nil
	}
	table := tablewriter.NewWriter(logrus.StandardLogger().Out)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.AppendBulk(data)
	table.Render()
	return nil
}

func (m *SMaintenance) List(svcProxyID string) (*[]models.Maintenance, error) {
	headers := httpclient.GetHeaders(m.Settings.SessionToken, m.Settings.Version, m.Settings.Pod, m.Settings.UsersID)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/maintenance", m.Settings.PaasHost, m.Settings.PaasHostVersion, m.Settings.EnvironmentID, svcProxyID), headers)
	if err != nil {
		return nil, err
	}
	var maintenance []models.Maintenance
	err = httpclient.ConvertResp(resp, statusCode, &maintenance)
	if err != nil {
		return nil, err
	}
	return &maintenance, nil
}
