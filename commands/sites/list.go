package sites

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/forana/simpletable"
)

func CmdList(is ISites, iservices services.IServices) error {
	serviceProxy, err := iservices.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}
	sites, err := is.List(serviceProxy.ID)
	if err != nil {
		return err
	}
	if sites == nil || len(*sites) == 0 {
		logrus.Println("No sites found")
		return nil
	}
	table, err := simpletable.New(simpletable.HeadersForType(models.Site{}), *sites)
	if err != nil {
		return err
	}
	table.Write(logrus.StandardLogger().Out)
	return nil
}

func (s *SSites) List(svcID string) (*[]models.Site, error) {
	headers := httpclient.GetHeaders(s.Settings.SessionToken, s.Settings.Version, s.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/sites", s.Settings.PaasHost, s.Settings.PaasHostVersion, s.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return nil, err
	}
	var sites []models.Site
	err = httpclient.ConvertResp(resp, statusCode, &sites)
	if err != nil {
		return nil, err
	}
	return &sites, nil
}
