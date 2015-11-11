package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
)

// ListVars lists all environment variables.
func ListVars(settings *models.Settings) {
	helpers.SignIn(settings)
	envVars := helpers.ListEnvVars(settings)
	var keys []string
	for k := range envVars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("%s=%s\n", key, envVars[key])
	}
}

// SetVar adds a new environment variables or updates the value of an existing
// environment variables. Any changes to environment variables will not take
// effect until the service is redeployed by pushing new code or via
// `catalyze redeploy`.
func SetVar(variables []string, settings *models.Settings) {
	helpers.SignIn(settings)
	envVars := helpers.ListEnvVars(settings)

	envVarsMap := make(map[string]string, len(variables))
	for _, envVar := range variables {
		pieces := strings.SplitN(envVar, "=", 2)
		if len(pieces) != 2 {
			fmt.Printf("Invalid variable format. Expected <key>=<value> but got %s\n", envVar)
			os.Exit(1)
		}
		envVarsMap[pieces[0]] = pieces[1]
	}

	for key := range envVarsMap {
		if _, ok := envVars[key]; ok {
			helpers.UnsetEnvVar(key, settings)
		}
	}

	helpers.SetEnvVars(envVarsMap, settings)
	fmt.Println("Set. For these environment variables to take effect, you will need to redeploy your service with \"catalyze redeploy\"")
}

// UnsetVar deletes an environment variable. Any changes to environment variables
// will not take effect until the service is redeployed by pushing new code
// or via `catalyze redeploy`.
func UnsetVar(variable string, settings *models.Settings) {
	helpers.SignIn(settings)
	helpers.UnsetEnvVar(variable, settings)
	fmt.Println("Unset.")
}
