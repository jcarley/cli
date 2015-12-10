package commands

import (
	"fmt"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
)

// ListServices lists the names of all services for an environment.
func ListServices(settings *models.Settings) {
	helpers.SignIn(settings)
	env := helpers.RetrieveEnvironment("pod", settings)
	fmt.Println("NAME")
	for _, s := range *env.Data.Services {
		fmt.Printf("%s\n", s.Label)
	}
}
