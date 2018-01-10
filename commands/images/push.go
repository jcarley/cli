package images

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/lib/images"
	"github.com/daticahealth/cli/lib/prompts"
	"github.com/daticahealth/cli/models"
)

func cmdImagePush(envID, name string, user *models.User, ie environments.IEnvironments, ii images.IImages, ip prompts.IPrompts) error {
	env, err := ie.Retrieve(envID)
	if err != nil {
		return err
	}

	image, err := ii.Push(name, user, env, ip)
	if err != nil {
		return err
	}
	fullImageName := fmt.Sprintf("%s:%s", image.Name, image.Tag)
	logrus.Printf("\nSuccessfully pushed image %s\n", fullImageName)

	repo := ii.GetNotaryRepository(env.Pod, image.Name, user)
	rootKeyPath := ""

	if err = ii.CheckChangelist(repo, ip); err != nil {
		return nil
	}

	if _, err = ii.ListTargets(repo); err != nil {
		if !strings.Contains(err.Error(), images.MissingTrustData) {
			return err
		}
		logrus.Println("Initializing trust repository")
		if err = ii.InitNotaryRepo(repo, rootKeyPath); err != nil {
			return err
		}
		logrus.Printf("Initialized trust repository for %s\n", image.Name)
	}

	logrus.Printf(`Adding target "%s" to trust repository`, image.Tag)
	//TODO: better error printing here? Gets output twice if missing the proper keys to sign
	if err = ii.AddTargetHash(repo, image.Digest, image.Tag, true); err != nil {
		return err
	}
	logrus.Printf("\nSuccessfully signed image %s\n", fullImageName)
	return nil
}
