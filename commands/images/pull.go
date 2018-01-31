package images

import (
	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/lib/images"
	"github.com/daticahealth/cli/models"
)

func cmdImagePull(envID, name string, user *models.User, ie environments.IEnvironments, ii images.IImages) error {
	env, err := ie.Retrieve(envID)
	if err != nil {
		return err
	}
	repositoryName, tag, err := ii.GetGloballyUniqueNamespace(name, env, true)
	if err != nil {
		return err
	}
	if tag == "" {
		logrus.Printf("No tag specified. Using default tag '%s'\n", images.DefaultTag)
		tag = images.DefaultTag
	}
	logrus.Println("Verifying image has been signed...")
	repo := ii.GetNotaryRepository(env.Pod, repositoryName, user)
	target, err := ii.LookupTarget(repo, tag)
	if err != nil {
		logrus.Warnf("Content verification failed: %s\n", err.Error())
		return nil
	}
	return ii.Pull(repositoryName, target, user, env)
}
