package associated

import (
	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/models"
)

func CmdAssociated(ia IAssociated) error {
	envs, err := ia.Associated()
	if err != nil {
		return err
	}
	for envAlias, env := range envs {
		logrus.Printf(`%s:
    Environment ID:   %s
    Environment Name: %s
    Service ID:       %s
    Associated at:    %s
    Pod:              %s
    Organization ID:  %s
`, envAlias, env.EnvironmentID, env.Name, env.ServiceID, env.Directory, env.Pod, env.OrgID)
	}
	if len(envs) == 0 {
		logrus.Println("No environments have been associated")
	}
	return nil
}

// Associated lists all currently associated environments.
func (a *SAssociated) Associated() (map[string]models.AssociatedEnv, error) {
	return a.Settings.Environments, nil
}
