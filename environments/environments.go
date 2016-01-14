package environments

import (
	"fmt"

	"github.com/catalyzeio/cli/models"
)

// CmdEnvironments lists all environments which the user has access to
func CmdEnvironments(environments IEnvironments) error {
	envs, err := environments.List()
	if err != nil {
		return err
	}
	for _, env := range *envs {
		fmt.Printf("%+v", env)
		//fmt.Printf("%s: %s\n", env.Data.Name, env.ID)
	}
	if len(*envs) == 0 {
		fmt.Println("no environments found")
	}
	return nil
}

func (e *SEnvironments) List() (*[]models.Environment, error) {
	return nil, nil
}

func (e *SEnvironments) Retrieve() (*models.Environment, error) {
	return nil, nil
}
