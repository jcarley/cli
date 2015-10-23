package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/catalyzeio/catalyze/config"
	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
)

// Associate an environment so that commands can be run against it. This command
// no longer adds a git remote. See commands.AddRemote().
func Associate(envLabel string, serviceLabel string, alias string, remote string, defaultEnv bool, settings *models.Settings) {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		fmt.Println("No git repo found in the current directory")
		os.Exit(1)
	}
	helpers.SignIn(settings)
	fmt.Printf("Existing git remotes named \"%s\" will be overwritten\n", remote)
	envs := helpers.ListEnvironments("pod", settings)
	for _, env := range *envs {
		if env.Data.Name == envLabel {
			if env.State == "defined" {
				fmt.Printf("Your environment is not yet provisioned. Please visit https://dashboard.catalyze.io/environments/update/%s to finish provisioning your environment\n", env.ID)
				return
			}
			// would be nice to have some sort of global filter() function
			var chosenService models.Service
			if serviceLabel != "" {
				labels := []string{}
				for _, service := range *env.Data.Services {
					if service.Type == "code" {
						labels = append(labels, service.Label)
						if service.Label == serviceLabel {
							chosenService = service
							break
						}
					}
				}
				if chosenService.Type == "" {
					fmt.Printf("No code service found with name '%s'. Code services found: %s\n", serviceLabel, strings.Join(labels, ", "))
					os.Exit(1)
				}
			} else {
				for _, service := range *env.Data.Services {
					if service.Type == "code" {
						chosenService = service
						break
					}
				}
				if chosenService.Type == "" {
					fmt.Printf("No code service found for \"%s\" environment (ID = %s)\n", envLabel, settings.EnvironmentID)
					os.Exit(1)
				}
			}
			for _, r := range helpers.ListGitRemote() {
				if r == remote {
					helpers.RemoveGitRemote(remote)
					break
				}
			}
			helpers.AddGitRemote(remote, chosenService.Source)
			fmt.Printf("\"%s\" remote added.\n", remote)
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				panic(err)
			}
			name := alias
			if name == "" {
				name = envLabel
			}
			settings.Environments[name] = models.AssociatedEnv{
				EnvironmentID: env.ID,
				ServiceID:     chosenService.ID,
				Directory:     dir,
				Name:          envLabel,
			}
			if defaultEnv {
				settings.Default = name
			}
			config.DropBreadcrumb(name, settings)
			config.SaveSettings(settings)
			if len(settings.Environments) > 1 && settings.Default == "" {
				fmt.Printf("You now have %d environments associated. Consider running \"catalyze default ENV_NAME\" to set a default\n", len(settings.Environments))
			}
			fmt.Printf("Your git repository \"%s\"  has been associated with code service \"%s\" and environment \"%s\"\n", remote, serviceLabel, name)
			return
		}
	}
	fmt.Printf("No environment with label \"%s\" found\n", envLabel)
	os.Exit(1)
}

// Associated lists all currently associated environments.
func Associated(settings *models.Settings) {
	for envAlias, env := range settings.Environments {
		fmt.Printf(`%s:
    Environment ID:   %s
    Environment Name: %s
    Service ID:       %s
    Associated at:    %s
    Default:          %v
`, envAlias, env.EnvironmentID, env.Name, env.ServiceID, env.Directory, settings.Default == envAlias)
	}
	if len(settings.Environments) == 0 {
		fmt.Println("No environments have been associated")
	}
}

// SetDefault sets the default environment. This environment must already be
// associated. Any commands run outside of a git directory will use the default
// environment for context.
func SetDefault(alias string, settings *models.Settings) {
	var found bool
	for name := range settings.Environments {
		if name == alias {
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("No environment with an alias of \"%s\" has been associated. Please run \"catalyze associate\" first\n", alias)
		os.Exit(1)
	}
	settings.Default = alias
	config.SaveSettings(settings)
	fmt.Printf("%s is now the default environment\n", alias)
}

// Disassociate removes an existing association with the environment. The
// `catalyze` remote on the local github repo will *NOT* be removed.
func Disassociate(alias string, settings *models.Settings) {
	// DeleteBreadcrumb removes the environment from the settings.Environments
	// array for you
	config.DeleteBreadcrumb(alias, settings)
	fmt.Printf("WARNING: Your existing git remote *has not* been removed.\n\n")
	fmt.Println("Association cleared.")
}
