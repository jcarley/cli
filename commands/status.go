package commands

import (
	"fmt"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
)

// Status prints out an environment healthcheck. The status of the environment
// and every service in the environment is printed out.
func Status(settings *models.Settings) {
	helpers.SignIn(settings)
	env := helpers.RetrieveEnvironment("pod", settings)
	fmt.Printf("%s (environment ID = %s):\n", env.Data.Name, env.ID)
	for _, service := range *env.Data.Services {
		if service.Type != "utility" {
			if service.Type == "code" {
				switch service.Size.(type) {
				case string:
					printLegacySizing(&service)
				default:
					printNewSizing(&service)
				}
			} else {
				switch service.Size.(type) {
				case string:
					sizeString := service.Size.(string)
					defer fmt.Printf("\t%s (size = %s, image = %s, status = %s) ID: %s\n", service.Label, sizeString, service.Name, service.DeployStatus, service.ID)
				default:
					serviceSize := service.Size.(map[string]interface{})
					defer fmt.Printf("\t%s (ram = %.0f, storage = %.0f, behavior = %s, type = %s, cpu = %.0f, image = %s, status = %s) ID: %s\n", service.Label, serviceSize["ram"], serviceSize["storage"], serviceSize["behavior"], serviceSize["type"], serviceSize["cpu"], service.Name, service.DeployStatus, service.ID)
				}
			}
		}
	}
}

func printLegacySizing(service *models.Service) {
	sizeString := service.Size.(string)
	fmt.Printf("\t%s (size = %s, build status = %s, deploy status = %s) ID: %s\n", service.Label, sizeString, service.BuildStatus, service.DeployStatus, service.ID)
}

func printNewSizing(service *models.Service) {
	serviceSize := service.Size.(map[string]interface{})
	fmt.Printf("\t%s (ram = %.0f, storage = %.0f, behavior = %s, type = %s, cpu = %.0f, build status = %s, deploy status = %s) ID: %s\n", service.Label, serviceSize["ram"], serviceSize["storage"], serviceSize["behavior"], serviceSize["type"], serviceSize["cpu"], service.BuildStatus, service.DeployStatus, service.ID)
}
