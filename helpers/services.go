package helpers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/catalyzeio/catalyze/httpclient"
	"github.com/catalyzeio/catalyze/models"
)

// RetrieveService returns a service model for the associated ServiceID
func RetrieveService(settings *models.Settings) *models.Service {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/services/%s", settings.PaasHost, settings.EnvironmentID, settings.ServiceID), true, settings)
	var service models.Service
	json.Unmarshal(resp, &service)
	return &service
}

// RetrieveServiceByLabel returns a Service object given its label by fetching
// the associated environment, looping through its services, and returning
// the first one matching the given label. If no service with the given label
// is found, nil is returned.
func RetrieveServiceByLabel(label string, settings *models.Settings) *models.Service {
	env := RetrieveEnvironment("pod", settings)
	for _, service := range *env.Data.Services {
		if service.Label == label {
			return &service
		}
	}
	return nil
}

// RetrieveTempURL fetches a temporary URL that can be used for downloading an
// existing database backup job. These URLs are signed and only valid for a
// short period of time.
func RetrieveTempURL(backupID string, serviceID string, settings *models.Settings) *models.TempURL {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/services/%s/backup/%s/url", settings.PaasHost, settings.EnvironmentID, serviceID, backupID), true, settings)
	var tempURL models.TempURL
	json.Unmarshal(resp, &tempURL)
	return &tempURL
}

// RetrieveTempLogsURL fetches a temporary URL that can be used for downloading
// logs for a 'finished' Backup/Restore/Import/Export job. These URLs are
// signed and only valid for a short period of time.
func RetrieveTempLogsURL(jobID string, jobType string, serviceID string, settings *models.Settings) *models.TempURL {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/services/%s/%s/%s/logs/url", settings.PaasHost, settings.EnvironmentID, serviceID, jobType, jobID), true, settings)
	var tempURL models.TempURL
	json.Unmarshal(resp, &tempURL)
	return &tempURL
}

// RetrieveTempUploadURL fetches a temporary URL that can be used for uploading
// a file for an Import job. These URLs are signed and only valid for a short
// period of time.
func RetrieveTempUploadURL(serviceID string, settings *models.Settings) *models.TempURL {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/services/%s/restore/url", settings.PaasHost, settings.EnvironmentID, serviceID), true, settings)
	var tempURL models.TempURL
	json.Unmarshal(resp, &tempURL)
	return &tempURL
}

// ListEnvVars returns all env vars for the associated Service
func ListEnvVars(settings *models.Settings) map[string]string {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/services/%s/env", settings.PaasHost, settings.EnvironmentID, settings.ServiceID), true, settings)
	var envVars map[string]string
	json.Unmarshal(resp, &envVars)
	return envVars
}

// ListBackups returns a list of all backups regardless of their status for a
// given service
func ListBackups(serviceID string, page int, pageSize int, settings *models.Settings) *[]models.Job {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/services/%s/backup?pageNum=%d&pageSize=%d", settings.PaasHost, settings.EnvironmentID, serviceID, page, pageSize), true, settings)
	var jobsMap map[string]models.Job
	json.Unmarshal(resp, &jobsMap)
	var jobs []models.Job
	for jobID, job := range jobsMap {
		job.ID = jobID
		jobs = append(jobs, job)
	}
	return &jobs
}

// CreateBackup kicks off a backup job for the given service. The task is
// returned which can be used to check on the status of the backup.
func CreateBackup(serviceID string, settings *models.Settings) *models.Task {
	backup := map[string]string{
		"archiveType":    "cf",
		"encryptionType": "aes",
	}
	b, err := json.Marshal(backup)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	resp := httpclient.Post(b, fmt.Sprintf("%s/v1/environments/%s/services/%s/backup", settings.PaasHost, settings.EnvironmentID, serviceID), true, settings)
	var m map[string]string
	json.Unmarshal(resp, &m)
	return &models.Task{
		ID: m["taskId"],
	}
}

// RestoreBackup kicks off a restore job for the given service. The task is
// returned which can be used to check on the status of the restore.
func RestoreBackup(serviceID string, backupID string, settings *models.Settings) *models.Task {
	backup := map[string]string{
		"archiveType":    "cf",
		"encryptionType": "aes",
	}
	b, err := json.Marshal(backup)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	resp := httpclient.Post(b, fmt.Sprintf("%s/v1/environments/%s/services/%s/restore/%s", settings.PaasHost, settings.EnvironmentID, serviceID, backupID), true, settings)
	var m map[string]string
	json.Unmarshal(resp, &m)
	return &models.Task{
		ID: m["taskId"],
	}
}

// InitiateRakeTask kicks off a rake task for the associated code service. The
// logs for the rake task are viewable in the environments logging server.
func InitiateRakeTask(taskName string, settings *models.Settings) {
	rakeTask := map[string]string{}
	b, err := json.Marshal(rakeTask)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	encodedTaskName, err := url.Parse(taskName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	httpclient.Post(b, fmt.Sprintf("%s/v1/environments/%s/services/%s/rake/%s", settings.PaasHost, settings.EnvironmentID, settings.ServiceID, encodedTaskName), true, settings)
}

// InitiateWorker starts a background worker for the associated code service
// for the given Procfile target.
func InitiateWorker(target string, settings *models.Settings) {
	worker := map[string]string{
		"target": target,
	}
	b, err := json.Marshal(worker)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	httpclient.Post(b, fmt.Sprintf("%s/v1/environments/%s/services/%s/background", settings.PaasHost, settings.EnvironmentID, settings.ServiceID), true, settings)
}

// RedeployService redeploys the associated code service
func RedeployService(serviceID string, settings *models.Settings) {
	redeploy := map[string]string{}
	b, err := json.Marshal(redeploy)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	httpclient.Post(b, fmt.Sprintf("%s/v1/environments/%s/services/%s/redeploy", settings.PaasHost, settings.EnvironmentID, serviceID), true, settings)
}

// InitiateImport starts an import job for the given database service
func InitiateImport(tempURL string, filePath string, key string, iv string, options map[string]string, wipeFirst bool, serviceID string, settings *models.Settings) *models.Task {
	httpclient.PutFile(filePath, tempURL, true, settings)
	importParams := models.Import{
		Location:  tempURL,
		Key:       key,
		IV:        iv,
		WipeFirst: wipeFirst,
		Options:   options,
	}
	b, err := json.Marshal(importParams)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	resp := httpclient.Post(b, fmt.Sprintf("%s/v1/environments/%s/services/%s/db/import", settings.PaasHost, settings.EnvironmentID, serviceID), true, settings)
	var task models.Task
	json.Unmarshal(resp, &task)
	return &task
}

// RequestConsole asks for a console to be setup. The console is not immediately
// ready but the resulting taskID should be used to check on its status.
func RequestConsole(command string, serviceID string, settings *models.Settings) *models.Task {
	console := map[string]string{}
	if command != "" {
		console["command"] = command
	}
	b, err := json.Marshal(console)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	resp := httpclient.Post(b, fmt.Sprintf("%s/v1/environments/%s/services/%s/console", settings.PaasHost, settings.EnvironmentID, serviceID), true, settings)
	var m map[string]string
	json.Unmarshal(resp, &m)
	return &models.Task{
		ID: m["taskId"],
	}
}

// RetrieveConsoleTokens returns the information necessary for connecting to a
// console service. The console service must already be ready and awaiting a
// connection.
func RetrieveConsoleTokens(jobID string, serviceID string, settings *models.Settings) *models.ConsoleCredentials {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/services/%s/console/token/%s", settings.PaasHost, settings.EnvironmentID, serviceID, jobID), true, settings)
	var credentials models.ConsoleCredentials
	json.Unmarshal(resp, &credentials)
	return &credentials
}

// DestroyConsole properly shuts down a console service
func DestroyConsole(jobID string, serviceID string, settings *models.Settings) {
	httpclient.Delete(fmt.Sprintf("%s/v1/environments/%s/services/%s/console/%s", settings.PaasHost, settings.EnvironmentID, serviceID, jobID), true, settings)
}

// RetrieveServiceMetrics fetches metrics for a single service for a specified
// number of minutes.
func RetrieveServiceMetrics(mins int, settings *models.Settings) *models.Metrics {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/metrics/%s?mins=%d", settings.PaasHost, settings.EnvironmentID, settings.ServiceID, mins), true, settings)
	var metrics models.Metrics
	json.Unmarshal(resp, &metrics)
	return &metrics
}

// ListServiceFiles retrieves a list of all downloadable service files for the
// specified code service.
func ListServiceFiles(serviceID string, settings *models.Settings) *[]models.ServiceFile {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/services/%s/files", settings.PaasHost, settings.EnvironmentID, serviceID), true, settings)
	var files []models.ServiceFile
	json.Unmarshal(resp, &files)
	return &files
}

// RetrieveServiceFile retrieves a service file by its ID.
func RetrieveServiceFile(serviceID, fileID string, settings *models.Settings) *models.ServiceFile {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/environments/%s/services/%s/files/%s", settings.PaasHost, settings.EnvironmentID, serviceID, fileID), true, settings)
	var file models.ServiceFile
	json.Unmarshal(resp, &file)
	return &file
}
