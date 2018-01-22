package targets

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/lib/images"
	"github.com/daticahealth/cli/models"
	"github.com/olekukonko/tablewriter"
)

func cmdTargetsList(envID, imageName string, user *models.User, ie environments.IEnvironments, ii images.IImages) error {
	env, err := ie.Retrieve(envID)
	if err != nil {
		return err
	}

	var targets []*images.Target

	repositoryName, tag, err := ii.GetGloballyUniqueNamespace(imageName, env)
	if err != nil {
		return err
	}
	repo := ii.GetNotaryRepository(env.Pod, repositoryName, user)
	if tag == "" {
		logrus.Printf("Searching for signed targets in trust repository %s\n", repositoryName)
		targets, err = ii.ListTargets(repo)
		if err != nil {
			return err
		}
	} else {
		logrus.Printf("Searching for signed target \"%s\" in trust repository %s\n", tag, repositoryName)
		target, err := ii.LookupTarget(repo, tag)
		if err != nil {
			return err
		}
		targets = append(targets, target)
	}

	if len(targets) > 0 {
		data := [][]string{{"Name", "Digest", "Size", "Role"}, {"----", "------", "----", "----"}}
		for _, t := range targets {
			data = append(data, []string{t.Name, string(t.Digest), fmt.Sprintf("%v", t.Size), t.Role})
		}

		table := tablewriter.NewWriter(logrus.StandardLogger().Out)
		table.SetBorder(false)
		table.SetRowLine(false)
		table.SetAlignment(1)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.AppendBulk(data)
		table.Render()
		logrus.Println()
	} else {
		logrus.Printf("No signed targets in remote trust repository %s\n", repositoryName)
	}
	return nil
}
