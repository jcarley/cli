package commands

import (
	"fmt"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
)

// Environments lists all environments which the user has access to
func Environments(settings *models.Settings) {
	helpers.SignIn(settings)
	envs := helpers.ListEnvironments("spec", settings)
	for _, env := range *envs {
		fmt.Printf("%s: %s\n", env.Data.Name, env.ID)
	}
	if len(*envs) == 0 {
		fmt.Println("no environments found")
	}
}
