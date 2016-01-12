package worker

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

// Worker starts a Procfile target as a worker. Worker containers are intended
// to be short-lived, one-off tasks.
func Worker(target string, settings *models.Settings) {
	helpers.SignIn(settings)
	fmt.Printf("Initiating a background worker for Service: %s (procfile target = \"%s\")\n", settings.ServiceID, target)
	helpers.InitiateWorker(target, settings)
	fmt.Println("Worker started.")
}
