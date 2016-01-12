package environments

import "github.com/catalyzeio/cli/models"

// IEnvironments is an interface for interacting with environments
type IEnvironments interface {
	List() (*[]models.Environment, error)
	Retrieve(envID string) (*models.Environment, error)
}

// SEnvironments is a concrete implementation of IEnvironments
type SEnvironments struct {
	Settings *models.Settings
}

// New generates a new instance of IEnvironments
func New(settings *models.Settings) IEnvironments {
	return &SEnvironments{
		Settings: settings,
	}
}
