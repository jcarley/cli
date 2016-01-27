package redeploy

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func CmdRedeploy(svcName string, ir IRedeploy, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the name \"%s\"", svcName)
	}
	logrus.Printf("Redeploying %s (ID = %s)", svcName, service.ID)
	err = ir.Redeploy(service)
	if err != nil {
		return err
	}
	logrus.Println("Redeploy successful! Check the status and logs for updates")
	return nil
}

// Redeploy offers a way of deploying a service without having to do a git push
// first. The same version of the currently running service is deployed with
// no changes.
func (r *SRedeploy) Redeploy(service *models.Service) error {
	headers := httpclient.GetHeaders(r.Settings.SessionToken, r.Settings.Version, r.Settings.Pod)
	resp, statusCode, err := httpclient.Post(nil, fmt.Sprintf("%s%s/services/%s/redeploy", r.Settings.PaasHost, r.Settings.PaasHostVersion, service.ID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
