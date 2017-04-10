package volumes

import "github.com/daticahealth/cli/models"

// IVolumes
type IVolumes interface {
	List(svcID string) (*[]models.Volume, error)
}

// SVolumes is a concrete implementation of IVolumes
type SVolumes struct {
	Settings *models.Settings
}

// New returns an instance of IVolumes
func New(settings *models.Settings) IVolumes {
	return &SVolumes{
		Settings: settings,
	}
}
