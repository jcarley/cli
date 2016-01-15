package redeploy

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
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
	helpers.RedeployService(service.ID, r.Settings)
	return nil
}
