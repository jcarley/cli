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
		ID:    "1234",
		State: "running",
		PodID: "pod01",
		Name:  "env",
		Pod:   "pod01",
		Services: &[]models.Service{{
			ID:    "1234",
			Type:  "code",
			Label: "app01",
			Size: map[string]string{
				"cpu": "1",
			},
			BuildStatus:  "finished",
			DeployStatus: "running",
			Name:         "svc",
			EnvVars:      map[string]string{},
			Source:       "git@git",
			LBIP:         "1.2.3.4",
			DockerImage:  "image",
		}},
		Namespace: "pod01",
		DNSName:   "pod01",
	}
}
