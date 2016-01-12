package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
)

// RetrieveJob fetches a Job model by its ID
func RetrieveJob(jobID string, serviceID string, settings *models.Settings) *models.Job {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/services/%s/jobs/%s", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID, jobID), true, settings)
	var job models.Job
	json.Unmarshal(resp, &job)
	return &job
}

// RetrieveJobFromTaskID translates a task into a job
func RetrieveJobFromTaskID(taskID string, settings *models.Settings) *models.Job {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/tasks/%s", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, taskID), true, settings)
	var job models.Job
	json.Unmarshal(resp, &job)
	return &job
}

// RetrieveRunningJobs fetches all running jobs for a service
func RetrieveRunningJobs(serviceID string, settings *models.Settings) *map[string]models.Job {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/services/%s/jobs?status=running", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID), true, settings)
	var jobs map[string]models.Job
	json.Unmarshal(resp, &jobs)
	return &jobs
}

// RetrieveLatestBuildJob fetches the latest build of a code service (nested in a json)
func RetrieveLatestBuildJob(serviceID string, settings *models.Settings) *map[string]models.Job {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/services/%s/jobs?type=build&pageSize=1", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID), true, settings)
	var jobs map[string]models.Job
	json.Unmarshal(resp, &jobs)
	return &jobs
}
