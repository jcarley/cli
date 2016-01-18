package associate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/environments"
	"github.com/catalyzeio/cli/git"
	"github.com/catalyzeio/cli/models"
)

func CmdAssociate(envLabel, svcLabel, alias, remote string, defaultEnv bool, ia IAssociate, ig git.IGit, ie environments.IEnvironments) error {
	if !ig.Exists() {
		return errors.New("No git repo found in the current directory")
	}
	logrus.Printf("Existing git remotes named \"%s\" will be overwritten", remote)
	envs, err := ie.List()
	if err != nil {
		return err
	}
	var e *models.Environment
	for _, env := range *envs {
		logrus.Debugf("associate env %+v", env)
		if env.Name == envLabel {
			pod := env.Pod
			e, err = ie.Retrieve(env.ID)
			if err != nil {
				return err
			}
			e.Pod = pod
			if e.State == "defined" {
				return fmt.Errorf("Your environment is not yet provisioned. Please visit https://dashboard.catalyze.io/environments/update/%s to finish provisioning your environment", env.ID)
			}
			break
		}
	}
	if e == nil {
		return fmt.Errorf("No environment with label \"%s\" found", envLabel)
	}

	var chosenService *models.Service
	availableCodeServices := []string{}
	for _, service := range *e.Services {
		if service.Type == "code" {
			if service.Label == svcLabel {
				chosenService = &service
				break
			}
			availableCodeServices = append(availableCodeServices, service.Label)
		}
	}
	if chosenService == nil {
		return fmt.Errorf("No code service found with name '%s'. Code services found: %s", svcLabel, strings.Join(availableCodeServices, ", "))
	}
	remotes, err := ig.List()
	if err != nil {
		return err
	}
	for _, r := range remotes {
		if r == remote {
			ig.Rm(remote)
			break
		}
	}
	err = ig.Add(remote, chosenService.Source)
	if err != nil {
		return err
	}
	logrus.Printf("\"%s\" remote added.", remote)

	name := alias
	if name == "" {
		name = envLabel
	}
	err = ia.Associate(name, remote, defaultEnv, e, chosenService)
	if err != nil {
		return err
	}
	logrus.Printf("Your git repository \"%s\"  has been associated with code service \"%s\" and environment \"%s\"", remote, svcLabel, name)
	return nil
}

// Associate an environment so that commands can be run against it. This command
// no longer adds a git remote. See commands.AddRemote().
func (s *SAssociate) Associate(name, remote string, defaultEnv bool, env *models.Environment, chosenService *models.Service) error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	s.Settings.Environments[name] = models.AssociatedEnv{
		EnvironmentID: env.ID,
		ServiceID:     chosenService.ID,
		Directory:     dir,
		Name:          env.Name,
		Pod:           env.Pod,
	}
	if defaultEnv {
		s.Settings.Default = name
	}
	config.DropBreadcrumb(name, s.Settings)
	if len(s.Settings.Environments) > 1 && s.Settings.Default == "" {
		logrus.Printf("You now have %d environments associated. Consider running \"catalyze default ENV_NAME\" to set a default", len(s.Settings.Environments))
	}

	return nil
}
