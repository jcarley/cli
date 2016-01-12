package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/catalyzeio/cli/models"
	"github.com/mitchellh/go-homedir"
)

// SettingsPath is the location of the catalyze config file.
const SettingsFile = ".catalyze"

// LocalSettingsPath stores a breadcrumb in a local git repo with an env name
const LocalSettingsFile = "catalyze-config.json"

// SettingsRetriever defines an interface for a class responsible for generating
// a settings object used for most commands in the CLI. Some examples might be
// for retrieving settings based on the settings file or generating a settings
// object based on a directly entered environment ID and service ID.
type SettingsRetriever interface {
	GetSettings(bool, bool, string, string, string, string, string, string) *models.Settings
}

// FileSettingsRetriever reads in data from the SettingsFile and generates a
// settings object.
type FileSettingsRetriever struct{}

// GetSettings returns a Settings object for the current context
func (s FileSettingsRetriever) GetSettings(required bool, promptForEnv bool, envName string, svcName string, baasHost string, paasHost string, username string, password string) *models.Settings {
	HomeDir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	file, err := os.Open(filepath.Join(HomeDir, SettingsFile))
	if os.IsNotExist(err) {
		file, err = os.Create(filepath.Join(HomeDir, SettingsFile))
	}
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var settings models.Settings
	json.NewDecoder(file).Decode(&settings)
	// would be best to default this to an initialized map rather than nil
	if settings.Environments == nil {
		settings.Environments = make(map[string]models.AssociatedEnv)
	}

	// try and set the given env first, if it exists
	if envName != "" {
		setGivenEnv(envName, &settings)
		if settings.EnvironmentID == "" || settings.ServiceID == "" {
			fmt.Printf("No environment named \"%s\" has been associated. Run \"catalyze associated\" to see what environments have been associated or run \"catalyze associate\" from a local git repo to create a new association\n", envName)
			os.Exit(1)
		}
	}

	// if no env name was given, try and fetch the local env
	if settings.EnvironmentID == "" || settings.ServiceID == "" {
		setLocalEnv(required, &settings)
	}

	// if its not there, fetch the default
	if settings.EnvironmentID == "" || settings.ServiceID == "" {
		setDefaultEnv(&settings)
	}

	// if no default, fetch the first associated env and print warning
	if settings.EnvironmentID == "" || settings.ServiceID == "" {
		// warn and ask
		setFirstAssociatedEnv(required, promptForEnv, &settings)
	}

	// if no env found, warn and quit
	if required && (settings.EnvironmentID == "" || settings.ServiceID == "") {
		fmt.Println("No Catalyze environment has been associated. Run \"catalyze associate\" from a local git repo first")
		os.Exit(1)
	}

	settings.BaasHost = baasHost
	settings.PaasHost = paasHost
	settings.Username = username
	settings.Password = password

	return &settings
}

// SaveSettings persists the settings to disk
func SaveSettings(settings *models.Settings) {
	HomeDir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	b, _ := json.Marshal(&settings)
	err = ioutil.WriteFile(filepath.Join(HomeDir, SettingsFile), b, 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// DropBreadcrumb drops a config file in a local git repo with the name of an
// environment and adds it to the list of global associated environments.
func DropBreadcrumb(envName string, settings *models.Settings) {
	b, _ := json.Marshal(&models.Breadcrumb{
		EnvName: envName,
	})
	err := ioutil.WriteFile(filepath.Join(".git", LocalSettingsFile), b, 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// DeleteBreadcrumb removes the config file at LocalSettingsPath
func DeleteBreadcrumb(alias string, settings *models.Settings) {
	env := settings.Environments[alias]
	dir := env.Directory
	dir = filepath.Join(dir, ".git", LocalSettingsFile)
	defer os.Remove(dir)

	delete(settings.Environments, alias)
	if settings.Default == alias {
		settings.Default = ""
	}
	os.Remove(dir)
	SaveSettings(settings)
}

// setGivenEnv takes the given env name and finds it in the env list
// in the given settings object. It then populates the EnvironmentID and
// ServiceID on the settings object with appropriate values.
func setGivenEnv(envName string, settings *models.Settings) {
	for eName, e := range settings.Environments {
		if eName == envName {
			settings.EnvironmentID = e.EnvironmentID
			settings.ServiceID = e.ServiceID
			settings.Pod = e.Pod
			settings.EnvironmentName = envName
		}
	}
}

// setLocalEnv searches .git/catalyze-config.json for an associated env and
// searches for it in the given settings object. It then populates the
// EnvironmentID and ServiceID on the settings object with appropriate values.
func setLocalEnv(required bool, settings *models.Settings) {
	file, err := os.Open(filepath.Join(".git", LocalSettingsFile))
	defer file.Close()
	if err == nil {
		var breadcrumb models.Breadcrumb
		json.NewDecoder(file).Decode(&breadcrumb)
		if breadcrumb.EnvironmentID != "" && required {
			// we found an old config file, try and translate it
			//convertSettings(&breadcrumb, settings)
			// or punt
			fmt.Println("Please reassociate your environment and then run this command again")
			os.Exit(1)
		}
		setGivenEnv(breadcrumb.EnvName, settings)
	}
}

// setDefaultEnv takes the name of the default env (if it exists) and finds it
// in the env list in the given settings object. It then populates the
// EnvironmentID and ServiceID on the settings object with appropriate values.
func setDefaultEnv(settings *models.Settings) {
	setGivenEnv(settings.Default, settings)
}

// setFirstAssociatedEnv is the last line of defense. If no other environments
// were found locally or from the default flag, then the first one in the list
// of environments in the given settings object is used to populate
// EnvironmentID and ServiceID with appropriate values.
func setFirstAssociatedEnv(required bool, promptForEnv bool, settings *models.Settings) {
	for envName, e := range settings.Environments {
		settings.EnvironmentID = e.EnvironmentID
		settings.ServiceID = e.ServiceID
		settings.Pod = e.Pod
		settings.EnvironmentName = e.Name
		if promptForEnv {
			defaultEnvPrompt(envName)
		}
		break
	}
}

// defaultEnvPrompt asks the user when they dont have a default environment and
// aren't in an associated directory if they would like to proceed with the
// first environment found.
func defaultEnvPrompt(envName string) {
	var answer string
	for {
		fmt.Printf("No environment was specified and no default environment was found. Falling back to %s\n", envName)
		fmt.Print("Do you wish to proceed? (y/n) ")
		fmt.Scanln(&answer)
		fmt.Println("")
		if answer != "y" && answer != "n" {
			fmt.Printf("%s is not a valid option. Please enter 'y' or 'n'\n", answer)
		} else {
			break
		}
	}
	if answer == "n" {
		fmt.Println("Exiting")
		os.Exit(1)
	}
}
