package targets

import (
	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/lib/images"
	"github.com/daticahealth/cli/models"
)

const (
	improperImageName = "Improperly formatted image name. Follow the convention <image>:<tag>"
)

func cmdTargetsReset(envID, imageName string, user *models.User, ie environments.IEnvironments, ii images.IImages) error {
	env, err := ie.Retrieve(envID)
	if err != nil {
		return err
	}

	repositoryName, tag, err := ii.GetGloballyUniqueNamespace(imageName, env, true)
	if err != nil {
		return err
	}
	repo := ii.GetNotaryRepository(env.Pod, repositoryName, user)

	changelist, err := repo.GetChangelist()
	if err != nil {
		return err
	}
	changes := changelist.List()
	if len(changes) > 0 {
		if tag != "" {
			var indices []int
			for i, change := range changes {
				if change.Path() == tag {
					indices = append(indices, i)
				}
			}
			err := changelist.Remove(indices)
			if err != nil {
				return err
			}
			logrus.Printf("Local changelist cleared for target \"%s\" in trust repository %s", tag, repositoryName)
			return nil
		} else {
			err := changelist.Clear("")
			if err != nil {
				return err
			}
		}
		logrus.Printf("Local changelist cleared for trust repository %s", repositoryName)
		return nil
	}
	logrus.Printf("No unpublished changes for trust repository %s\n", repositoryName)
	return nil
}
