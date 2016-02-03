package environments

import (
	"errors"

	"github.com/catalyzeio/cli/models"
)

// MEnvironments is a mock implementation of IEnvironments
type MEnvironments struct {
	ReturnError bool
}

func (e *MEnvironments) List() (*[]models.Environment, error) {
	if e.ReturnError {
		return nil, errors.New("Mock error returned")
	}
	return &[]models.Environment{env()}, nil
}

func (e *MEnvironments) Retrieve(envID string) (*models.Environment, error) {
	if e.ReturnError {
		return nil, errors.New("Mock error returned")
	}
	environment := env()
	return &environment, nil
}

func env() models.Environment {
	return models.Environment{
		ID:        "1234",
		Name:      "env",
		Pod:       "pod01",
		Namespace: "pod01",
		DNSName:   "pod01",
	}
}
