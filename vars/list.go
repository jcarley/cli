package vars

import (
	"fmt"
	"sort"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
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
