package releases

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/olekukonko/tablewriter"
)

func CmdList(svcName string, ir IReleases, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services\" command.", svcName)
	}

	rls, err := ir.List(service.ID)
	if err != nil {
		return err
	}

	if rls == nil || len(*rls) == 0 {
		logrus.Println("No releases found")
		return nil
	}

	data := [][]string{{"Release Name", "Created At", "Notes"}}
	for _, r := range *rls {
		data = append(data, []string{r.ID, r.CreatedAt, r.Notes})
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

func (r *SReleases) List(svcID string) (*[]models.Release, error) {
	headers := httpclient.GetHeaders(r.Settings.SessionToken, r.Settings.Version, r.Settings.Pod, r.Settings.UsersID)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/releases", r.Settings.PaasHost, r.Settings.PaasHostVersion, r.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return nil, err
	}
	var rls []models.Release
	err = httpclient.ConvertResp(resp, statusCode, &rls)
	if err != nil {
		return nil, err
	}
	return &rls, nil
}
