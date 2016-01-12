package vars

import (
	"fmt"
	"os"
	"strings"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

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
