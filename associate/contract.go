package associate

import (
	"github.com/catalyzeio/cli/environments"
	"github.com/catalyzeio/cli/git"
	"github.com/catalyzeio/cli/models"
)

// interfaces are the API calls
type IAssociate interface {
	Associate() error
}

// SAssociate is a concrete implementation of IAssociate
type SAssociate struct {
	Settings     *models.Settings
	Git          git.IGit
	Environments environments.IEnvironments
	EnvLabel     string
	SvcLabel     string
	Alias        string
	Remote       string
	DefaultEnv   bool
}

func New(settings *models.Settings, git git.IGit, environments environments.IEnvironments, envLabel, svcLabel, alias, remote string, defaultEnv bool) IAssociate {
	return &SAssociate{
		Settings:     settings,
		Git:          git,
		Environments: environments,
		EnvLabel:     envLabel,
		SvcLabel:     svcLabel,
		Alias:        alias,
		Remote:       remote,
		DefaultEnv:   defaultEnv,
	}
}
