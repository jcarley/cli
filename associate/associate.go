package associate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/environments"
	"github.com/catalyzeio/cli/git"
	"github.com/catalyzeio/cli/models"
)

func CmdAssociate(ia IAssociate, ig git.IGit, ie environments.IEnvironments) error {
	if !s.Git.Exists() {
		return errors.New("No git repo found in the current directory")
	}
	// TODO fmt.Printf("Existing git remotes named \"%s\" will be overwritten\n", s.Remote)
	fmt.Printf("Existing git remotes will be overwritten\n", s.Remote)
	envs, err := ie.List()
	if err != nil {
		return err
	}
	var e *models.Environment
	for _, env := range *envs {
		fmt.Printf("\n\nassociate env %+v\n\n", env)
		if env.Name == s.EnvLabel {
			pod := env.Pod
			e, err = ie.Retrieve(env.ID)
			if err != nil {
				return err
			}
			e.Pod = pod
			if e.State == "defined" {
				return fmt.Errorf("Your environment is not yet provisioned. Please visit https://dashboard.catalyze.io/environments/update/%s to finish provisioning your environment\n", env.ID)
			}
			break
		}
	}
	if e == nil {
		// TODO return fmt.Errorf("No environment with label \"%s\" found\n", s.EnvLabel)
		return fmt.Errorf("No environment with given name was found")
	}

	var chosenService models.Service
	availableCodeServices := []string{}
	for _, service := range *e.Services {
		if service.Type == "code" {
			if service.Label == s.SvcLabel {
				chosenService = service
				break
			}
			availableCodeServices = append(availableCodeServices, service.Label)
		}
	}
	if chosenService.Type == "" {
		return fmt.Errorf("No code service found with name '%s'. Code services found: %s\n", s.SvcLabel, strings.Join(availableCodeServices, ", "))
	}
	remotes, err := s.Git.List()
	if err != nil {
		return err
	}
	for _, r := range remotes {
		if r == s.Remote {
			s.Git.Rm(s.Remote)
			break
		}
	}
	err = s.Git.Add(s.Remote, chosenService.Source)
	if err != nil {
		return err
	}
	fmt.Printf("\"%s\" remote added.\n", s.Remote)

}

// Associate an environment so that commands can be run against it. This command
// no longer adds a git remote. See commands.AddRemote().
func (s *SAssociate) Associate() error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	name := s.Alias
	if name == "" {
		name = s.EnvLabel
	}
	s.Settings.Environments[name] = models.AssociatedEnv{
		EnvironmentID: env.ID,
		ServiceID:     chosenService.ID,
		Directory:     dir,
		Name:          s.EnvLabel,
		Pod:           env.Pod,
	}
	if s.DefaultEnv {
		s.Settings.Default = name
	}
	config.DropBreadcrumb(name, s.Settings)
	config.SaveSettings(s.Settings)
	if len(s.Settings.Environments) > 1 && s.Settings.Default == "" {
		fmt.Printf("You now have %d environments associated. Consider running \"catalyze default ENV_NAME\" to set a default\n", len(s.Settings.Environments))
	}
	fmt.Printf("Your git repository \"%s\"  has been associated with code service \"%s\" and environment \"%s\"\n", s.Remote, s.SvcLabel, name)
	return nil
}
