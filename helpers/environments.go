package helpers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/catalyzeio/catalyze/config"
	"github.com/catalyzeio/catalyze/httpclient"
	"github.com/catalyzeio/catalyze/models"
)

// ListEnvironments returns a list of all environments the authorized
// user has access to
func ListEnvironments(source string, settings *models.Settings) *[]models.Environment {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments?pageSize=1000&source=%s", settings.PaasHost, config.PaasHostVersion, source), true, settings)
	var envs []models.Environment
	json.Unmarshal(resp, &envs)
	return &envs
}

// ListEnvironmentUsers returns a list of all users who have access to the
// associated environment
func ListEnvironmentUsers(settings *models.Settings) *models.EnvironmentUsers {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/users", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID), true, settings)
	var users models.EnvironmentUsers
	json.Unmarshal(resp, &users)
	return &users
}

// RetrieveEnvironment returns the associated environment model. The source
// parameter specifies where the data should come from. If source is `pod` then
// the Environment data will be fetched from the Pod API, otherwise the data
// will be retrieved from the Customer API.
func RetrieveEnvironment(source string, settings *models.Settings) *models.Environment {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s?source=%s", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, source), true, settings)
	var env models.Environment
	json.Unmarshal(resp, &env)
	return &env
}

// AddUserToEnvironment grants a user access to the associated env
func AddUserToEnvironment(usersID string, settings *models.Settings) {
	body := map[string]string{}
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/users/%s", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, usersID), true, settings)
}

// RemoveUserFromEnvironment revokes a users access to the associated env
func RemoveUserFromEnvironment(usersID string, settings *models.Settings) {
	httpclient.Delete(fmt.Sprintf("%s%s/environments/%s/users/%s", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, usersID), true, settings)
}

// SetEnvVars adds new env vars or updates the values of existing ones
func SetEnvVars(envVars map[string]string, settings *models.Settings) {
	b, err := json.Marshal(envVars)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/env", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, settings.ServiceID), true, settings)
}

// UnsetEnvVar deletes an env var from the associated code service
func UnsetEnvVar(envVar string, settings *models.Settings) {
	httpclient.Delete(fmt.Sprintf("%s%s/environments/%s/services/%s/env/%s", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, settings.ServiceID, envVar), true, settings)
}

// RetrieveEnvironmentMetrics fetches metrics for an entire environment for a
// specified number of minutes.
func RetrieveEnvironmentMetrics(mins int, settings *models.Settings) *[]models.Metrics {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/metrics?mins=%d", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, mins), true, settings)
	var metrics []models.Metrics
	json.Unmarshal(resp, &metrics)
	return &metrics
}
