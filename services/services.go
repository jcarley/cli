package services

import (
	"fmt"

	"github.com/catalyzeio/cli/models"
)

// CmdServices lists the names of all services for an environment.
func CmdServices(is IServices) error {
	svcs, err := is.List()
	if err != nil {
		return err
	}
	fmt.Println("NAME")
	for _, s := range *svcs {
		fmt.Printf("%s\n", s.Label)
	}
	return nil
}

func (s *SServices) List() (*[]models.Service, error) {
	return nil, nil
}

func (s *SServices) Retrieve(svcID string) (*models.Service, error) {
	return nil, nil
}

func (s *SServices) RetrieveByLabel(label string) (*models.Service, error) {
	return nil, nil
}
