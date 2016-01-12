package associate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/models"
)

// Associate an environment so that commands can be run against it. This command
// no longer adds a git remote. See commands.AddRemote().
func (s *SAssociate) Associate() error {
	if s.Git.Exists() {
		return errors.New("No git repo found in the current directory")
	}
	fmt.Printf("Existing git remotes named \"%s\" will be overwritten\n", s.Remote)
	envs, err := s.Environments.List()
	if err != nil {
		return err
	}
	for _, env := range *envs {
		fmt.Printf("\n\nassociate env %+v\n\n", env)
		if env.Name == s.EnvLabel {
			pod := env.Pod
			e, err := s.Environments.Retrieve(env.ID)
			if err != nil {
				return err
			}
			env = *e
			env.Pod = pod
			if env.State == "defined" {
				return fmt.Errorf("Your environment is not yet provisioned. Please visit https://dashboard.catalyze.io/environments/update/%s to finish provisioning your environment\n", env.ID)
			}

			var chosenService models.Service
			availableCodeServices := []string{}
			for _, service := range *env.Services {
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

			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				// TODO this is not a nice error, fix it
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
	}
	return fmt.Errorf("No environment with label \"%s\" found\n", s.EnvLabel)
}
