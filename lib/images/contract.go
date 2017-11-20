package images

import "github.com/daticahealth/cli/models"

// IImages describes container-image-related functionality
type IImages interface {
	ListImages() (*[]string, error)
	ListTags(imageName string) (*[]string, error)
	DeleteTag(imageName, tagName string) error
}

// SImages is a concrete implementation of IImages
type SImages struct {
	Settings *models.Settings
}

// New constructs an implementation of IImages
func New(settings *models.Settings) IImages {
	return &SImages{Settings: settings}
}
