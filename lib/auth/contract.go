package auth

import (
	"github.com/daticahealth/cli/lib/prompts"
	"github.com/daticahealth/cli/models"
)

// IAuth represents the contract that concrete implementations should follow
// when implementing authentication.
type IAuth interface {
	Signin() (*models.User, error)
	Signout() error
	Verify() (*models.User, error)
}

// SAuth is a concrete implementation of IAuth
type SAuth struct {
	Settings *models.Settings
	Prompts  prompts.IPrompts
}

func New(settings *models.Settings, prompts prompts.IPrompts) IAuth {
	return &SAuth{
		Settings: settings,
		Prompts:  prompts,
	}
}
