package environments

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

// Environments lists all environments which the user has access to
func Environments(settings *models.Settings) {
	helpers.SignIn(settings)
	envs := helpers.ListEnvironments(settings)
	for _, env := range *envs {
		fmt.Printf("%+v", env)
		//fmt.Printf("%s: %s\n", env.Data.Name, env.ID)
	}
	if len(*envs) == 0 {
		fmt.Println("no environments found")
	}
}

func (e *SEnvironments) List() (*[]models.Environment, error) {
	return nil, nil
}

func (e *SEnvironments) Retrieve(envID string) (*models.Environment, error) {
	return nil, nil
}
