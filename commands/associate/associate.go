package associate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/commands/git"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/models"
)

func CmdAssociate(envLabel, svcLabel, alias, remote string, defaultEnv bool, ia IAssociate, ig git.IGit, ie environments.IEnvironments, is services.IServices) error {
	logrus.Warnln("The \"--default\" flag has been deprecated! It will be removed in a future version.")
	logrus.Warnln("Please specify \"-E\" on all commands instead of using the default.")
	if !ig.Exists() {
		return errors.New("No git repo found in the current directory")
	}
	logrus.Printf("Existing git remotes named \"%s\" will be overwritten", remote)
	envs, errs := ie.List()
	if errs != nil && len(errs) > 0 {
		for pod, err := range errs {
			logrus.Debugf("Failed to list environments for pod \"%s\": %s", pod, err)
		}
	}
	var e *models.Environment
	var svcs *[]models.Service
	var err error
	for _, env := range *envs {
		if env.Name == envLabel {
			e = &env
			svcs, err = is.ListByEnvID(env.ID, env.Pod)
			if err != nil {
				return err
			}
			break
		}
	}
	if e == nil {
		return fmt.Errorf("No environment with name \"%s\" found", envLabel)
	}
	if svcs == nil {
		return fmt.Errorf("No services found for environment with name \"%s\"", envLabel)
	}

	var chosenService *models.Service
	availableCodeServices := []string{}
	for _, service := range *svcs {
		if service.Type == "code" {
			if service.Label == svcLabel {
				chosenService = &service
				break
			}
			availableCodeServices = append(availableCodeServices, service.Label)
		}
	}
	if chosenService == nil {
		return fmt.Errorf("No code service found with label \"%s\". Code services found: %s", svcLabel, strings.Join(availableCodeServices, ", "))
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
	logrus.Printf("Your git repository \"%s\" has been associated with code service \"%s\" and environment \"%s\"", remote, svcLabel, name)
	logrus.Println("After associating to an environment, you need to add a cert with the \"catalyze certs create\" command, if you have not done so already")
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
		OrgID:         env.OrgID,
	}

	return nil
}
