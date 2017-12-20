package targets

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/lib/images"
	"github.com/daticahealth/cli/lib/prompts"
	"github.com/daticahealth/cli/models"
)

func cmdTargetsDelete(envID, imageName string, user *models.User, ie environments.IEnvironments, ii images.IImages, ip prompts.IPrompts) error {
	env, err := ie.Retrieve(envID)
	if err != nil {
		return err
	}
	var targetList []string
	repositoryName, tag, err := ii.GetGloballyUniqueNamespace(imageName, env)
	if err != nil {
		return err
	}
	repo := ii.GetNotaryRepository(env.Pod, repositoryName, user)
	if tag == "" {
		if err := ip.YesNo("No tag specified", fmt.Sprintf("Would you like to delete trust data for all targets in %s? (y/n) ", repositoryName)); err != nil {
			return nil
		}

		targets, err := ii.ListTargets(repo)
		if err != nil {
			return err
		}
		for _, target := range targets {
			targetList = append(targetList, target.Name)
		}
	} else {
		targetList = []string{tag}
	}

	if err := ii.CheckChangelist(repo, ip); err != nil {
		return err
	}
	if err := ii.DeleteTargets(repo, targetList, true); err != nil {
		return err
	}
	logrus.Printf("\nSuccessfully deleted signed targets for %s\n", repositoryName)
	return nil
}
