package redeploy

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
)

func CmdRedeploy(svcName string, ir IRedeploy, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services\" command.", svcName)
	}
	logrus.Printf("Redeploying %s (ID = %s)", svcName, service.ID)
	err = ir.Redeploy("", service.ID)
	if err != nil {
		return err
	}
	logrus.Println("Redeploy successful! Check the status with \"catalyze status\" and your logging dashboard for updates")
	return nil
}

// Redeploy offers a way of deploying a service without having to do a git push
// first. The same version of the currently running service is deployed with
// no changes.
func (r *SRedeploy) Redeploy(releaseName, svcID string) error {
	var releaseParam = ""
	if releaseName != "" {
		releaseParam = fmt.Sprintf("&release=%s", releaseName)
	}
	headers := httpclient.GetHeaders(r.Settings.SessionToken, r.Settings.Version, r.Settings.Pod, r.Settings.UsersID)
	resp, statusCode, err := httpclient.Post(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/deploy?redeploy=true%s", r.Settings.PaasHost, r.Settings.PaasHostVersion, r.Settings.EnvironmentID, svcID, releaseParam), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
