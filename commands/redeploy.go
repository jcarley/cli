package commands

import (
	"fmt"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
)

// Redeploy offers a way of deploying a service without having to do a git push
// first. The same version of the currently running service is deployed with
// no changes.
func Redeploy(serviceLabel string, settings *models.Settings) {
	helpers.SignIn(settings)
	service := helpers.RetrieveServiceByLabel(serviceLabel, settings)
	if service == nil {
		panic(fmt.Errorf("Could not find a service with the name \"%s\"\n", serviceLabel))
	}
	fmt.Printf("Redeploying %s (ID = %s)\n", serviceLabel, service.ID)
	helpers.RedeployService(service.ID, settings)
	fmt.Println("Redeploy successful! Check the status and logs for updates")
}
