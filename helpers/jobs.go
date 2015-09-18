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

// RetrieveAllJobs fetches all running jobs for a service
func RetrieveAllJobs(serviceID string, settings *models.Settings) *map[string]models.Job {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/services/%s/jobs?pageSize=10000", settings.PaasHost, settings.EnvironmentID, serviceID), true, settings)
	var jobs map[string]models.Job
	json.Unmarshal(resp, &jobs)
	return &jobs
}
