package images

import (
	"github.com/daticahealth/cli/lib/prompts"
	"github.com/daticahealth/cli/models"
	notaryClient "github.com/docker/notary/client"
	"github.com/docker/notary/client/changelist"
)

// IImages describes container-image-related functionality
type IImages interface {
	ListImages() (*[]string, error)
	ListTags(imageName string) (*[]string, error)
	DeleteTag(imageName, tagName string) error
	Push(name string, user *models.User, env *models.Environment, ip prompts.IPrompts) (*models.Image, error)
	Pull(name string, user *models.User, env *models.Environment) (*models.Image, error)
	InitNotaryRepo(repo notaryClient.Repository, rootKeyPath string) error
	AddTargetHash(repo notaryClient.Repository, digest *models.ContentDigest, tag string, publish bool) error
	ListTargets(repo notaryClient.Repository, roles ...string) ([]*notaryClient.TargetWithRole, error)
	LookupTarget(repo notaryClient.Repository, tag string) (*notaryClient.TargetWithRole, error)
	DeleteTargets(repo notaryClient.Repository, tags []string, publish bool) error
	PrintChangelist(changes []changelist.Change)
	CheckChangelist(repo notaryClient.Repository, ip prompts.IPrompts) error
	GetNotaryRepository(pod, imageName string, user *models.User) notaryClient.Repository
	GetGloballyUniqueNamespace(name string, env *models.Environment) (string, string, error)
	Publish(repo notaryClient.Repository) error
}

// SImages is a concrete implementation of IImages
type SImages struct {
	Settings *models.Settings
}

// New constructs an implementation of IImages
func New(settings *models.Settings) IImages {
	return &SImages{Settings: settings}
}
