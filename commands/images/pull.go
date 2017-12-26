package images

import (
	"crypto/subtle"
	"encoding/hex"
	"fmt"

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

	image, err := ii.Pull(name, user, env)
	if err != nil {
		return err
	}
	fullImageName := fmt.Sprintf("%s:%s", image.Name, image.Tag)

	repo := ii.GetNotaryRepository(env.Pod, image.Name, user)
	target, err := ii.LookupTarget(repo, image.Tag)
	if err != nil {
		logrus.Warnf("Content verification failed. %s", err.Error())
		return nil
	}
	hashBytes, err := hex.DecodeString(image.Digest.Hash)
	if err != nil {
		return err
	}
	if subtle.ConstantTimeCompare(target.Hashes[image.Digest.HashType], hashBytes) == 0 {
		return fmt.Errorf("Content verification failed. Image content does not match the signed target.")
	}
	logrus.Printf("\nSuccessfully verified image %s against signed target in remote trust repository\n", fullImageName)

	return nil
}
