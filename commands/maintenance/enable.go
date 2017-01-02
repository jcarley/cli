package maintenance

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
)

func CmdEnable(svcName string, im IMaintenance, is services.IServices) error {
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

	err = im.Enable(serviceProxy.ID, upstreamService.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Maintenance mode enabled for service %s (ID = %s)", upstreamService.Label, upstreamService.ID)
	return nil
}

func (m *SMaintenance) Enable(svcProxyID, upstreamID string) error {
	body := map[string]string{
		"upstream": upstreamID,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(m.Settings.SessionToken, m.Settings.Version, m.Settings.Pod, m.Settings.UsersID)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/maintenance", m.Settings.PaasHost, m.Settings.PaasHostVersion, m.Settings.EnvironmentID, svcProxyID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
