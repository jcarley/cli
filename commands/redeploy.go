package commands

import (
	"fmt"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
)

// Redeploy offers a way of deploying a service without having to do a git push
// first. The same version of the currently running service is deployed with
// no changes.
func Redeploy(settings *models.Settings) {
	helpers.SignIn(settings)
	fmt.Printf("Redeploying %s\n", settings.ServiceID)
	helpers.RedeployService(settings)
	fmt.Println("Redeploy successful! Check the status and logs for updates")
}
