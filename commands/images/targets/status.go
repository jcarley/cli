package targets

import (
	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/lib/images"
	"github.com/daticahealth/cli/models"
	"github.com/docker/notary/client/changelist"
)

func cmdTargetsStatus(envID, imageName string, user *models.User, ie environments.IEnvironments, ii images.IImages) error {
	env, err := ie.Retrieve(envID)
	if err != nil {
		return err
	}

	repositoryName, tag, err := ii.GetGloballyUniqueNamespace(imageName, env)
	if err != nil {
		return err
	}
	repo := ii.GetNotaryRepository(env.Pod, repositoryName, user)
	changelist, err := repo.GetChangelist()
	if err != nil {
		return err
	}
	changes := filterChanges(changelist.List(), tag)
	if len(changes) > 0 {
		ii.PrintChangelist(changes)
	} else {
		logrus.Printf("No unpublished changes for %s\n", repositoryName)
	}
	return nil
}

func filterChanges(changes []changelist.Change, tag string) (filteredList []changelist.Change) {
	if tag == "" {
		return changes
	}
	for _, change := range changes {
		if change.Path() == tag {
			filteredList = append(filteredList, change)
		}
	}
	return filteredList
}
