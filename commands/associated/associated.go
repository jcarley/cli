package associated

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
)

func CmdAssociated(ia IAssociated) error {
	envs, defaultEnv, err := ia.Associated()
	if err != nil {
		return err
	}
	for envAlias, env := range envs {
		logrus.Printf(`%s:
    Environment ID:   %s
    Environment Name: %s
    Service ID:       %s
    Associated at:    %s
    Default:          %v
    Pod:              %s
`, envAlias, env.EnvironmentID, env.Name, env.ServiceID, env.Directory, defaultEnv == envAlias, env.Pod)
	}
	if len(envs) == 0 {
		logrus.Println("No environments have been associated")
	}
	return nil
}

// Associated lists all currently associated environments.
func (a *SAssociated) Associated() (map[string]models.AssociatedEnv, string, error) {
	return a.Settings.Environments, a.Settings.Default, nil
}
