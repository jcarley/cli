package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/catalyzeio/catalyze/httpclient"
	"github.com/catalyzeio/catalyze/models"
)

// RetrieveJob fetches a Job model by its ID
func RetrieveJob(jobID string, serviceID string, settings *models.Settings) *models.Job {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/services/%s/jobs/%s", settings.PaasHost, settings.EnvironmentID, serviceID, jobID), true, settings)
	var job models.Job
	json.Unmarshal(resp, &job)
	return &job
}

// RetrieveJobFromTaskID translates a task into a job
func RetrieveJobFromTaskID(taskID string, settings *models.Settings) *models.Job {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/tasks/%s", settings.PaasHost, settings.EnvironmentID, taskID), true, settings)
	var job models.Job
	json.Unmarshal(resp, &job)
	return &job
}

// RetrieveRunningJobs fetches all running jobs for a service
func RetrieveRunningJobs(serviceID string, settings *models.Settings) *map[string]models.Job {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/services/%s/jobs?status=running", settings.PaasHost, settings.EnvironmentID, serviceID), true, settings)
	var jobs map[string]models.Job
	json.Unmarshal(resp, &jobs)
	return &jobs
}

// RetrieveLatestBuildJob fetches the latest build of a code service (nested in a json)
func RetrieveLatestBuildJob(serviceID string, settings *models.Settings) *map[string]models.Job {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/services/%s/jobs?type=build&pageSize=1", settings.PaasHost, settings.EnvironmentID, serviceID), true, settings)
	var jobs map[string]models.Job
	json.Unmarshal(resp, &jobs)
	return &jobs
}
