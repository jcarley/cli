package services

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

// ListServices lists the names of all services for an environment.
func ListServices(settings *models.Settings) {
	helpers.SignIn(settings)
	env := helpers.RetrieveEnvironment("pod", settings)
	fmt.Println("NAME")
	for _, s := range *env.Services {
		fmt.Printf("%s\n", s.Label)
	}
}
