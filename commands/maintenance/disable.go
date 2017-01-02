package maintenance

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
)

func CmdDisable(svcName string, im IMaintenance, is services.IServices) error {
	// TODO need a way to determine if maintenance mode is available for an environment
	// preferrably without hardcoding
	upstreamService, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if upstreamService == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services\" command.", svcName)
	}

	serviceProxy, err := is.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}

	svcMaintenance, err := im.List(serviceProxy.ID)
	if err != nil {
		return err
	}
	enabled := false
	for _, mm := range *svcMaintenance {
		if mm.UpstreamID == upstreamService.ID {
			enabled = true
		}
	}
	if !enabled {
		return fmt.Errorf("Maintenance mode is not currently enabled for the service %s", svcName)
	}

	err = im.Disable(serviceProxy.ID, upstreamService.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Maintenance mode disabled for service %s (ID = %s)", upstreamService.Label, upstreamService.ID)
	return nil
}

func (m *SMaintenance) Disable(svcProxyID, upstreamID string) error {
	headers := httpclient.GetHeaders(m.Settings.SessionToken, m.Settings.Version, m.Settings.Pod, m.Settings.UsersID)
	resp, statusCode, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/maintenance?upstream=%s", m.Settings.PaasHost, m.Settings.PaasHostVersion, m.Settings.EnvironmentID, svcProxyID, upstreamID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
