package redeploy

import (
	"encoding/json"
	"fmt"

	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/services"
)

func CmdRedeploy(svcName string, ir IRedeploy, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the name \"%s\"\n", svcName)
	}
	fmt.Printf("Redeploying %s (ID = %s)\n", svcName, service.ID)
	err = ir.Redeploy(service)
	if err != nil {
		return err
	}
	fmt.Println("Redeploy successful! Check the status and logs for updates")
	return nil
}

// Redeploy offers a way of deploying a service without having to do a git push
// first. The same version of the currently running service is deployed with
// no changes.
func (r *SRedeploy) Redeploy(service *models.Service) error {
	redeploy := map[string]string{}
	b, err := json.Marshal(redeploy)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(r.Settings.APIKey, r.Settings.SessionToken, r.Settings.Version, r.Settings.Pod)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/redeploy", r.Settings.PaasHost, r.Settings.PaasHostVersion, r.Settings.EnvironmentID, service.ID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
