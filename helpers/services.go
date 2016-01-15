package helpers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
)

// RetrieveTempURL fetches a temporary URL that can be used for downloading an
// existing database backup job. These URLs are signed and only valid for a
// short period of time.
func RetrieveTempURL(backupID string, serviceID string, settings *models.Settings) *models.TempURL {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/services/%s/backup/%s/url", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID, backupID), true, settings)
	var tempURL models.TempURL
	json.Unmarshal(resp, &tempURL)
	return &tempURL
}

// RetrieveTempLogsURL fetches a temporary URL that can be used for downloading
// logs for a 'finished' Backup/Restore/Import/Export job. These URLs are
// signed and only valid for a short period of time.
func RetrieveTempLogsURL(jobID string, jobType string, serviceID string, settings *models.Settings) *models.TempURL {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/services/%s/%s/%s/logs/url", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID, jobType, jobID), true, settings)
	var tempURL models.TempURL
	json.Unmarshal(resp, &tempURL)
	return &tempURL
}

// RetrieveTempUploadURL fetches a temporary URL that can be used for uploading
// a file for an Import job. These URLs are signed and only valid for a short
// period of time.
func RetrieveTempUploadURL(serviceID string, settings *models.Settings) *models.TempURL {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/services/%s/restore/url", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID), true, settings)
	var tempURL models.TempURL
	json.Unmarshal(resp, &tempURL)
	return &tempURL
}

// ListBackups returns a list of all backups regardless of their status for a
// given service
func ListBackups(serviceID string, page int, pageSize int, settings *models.Settings) *[]models.Job {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/services/%s/backup?pageNum=%d&pageSize=%d", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID, page, pageSize), true, settings)
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
	resp := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/backup", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID), true, settings)
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
	resp := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/restore/%s", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID, backupID), true, settings)
	var m map[string]string
	json.Unmarshal(resp, &m)
	return &models.Task{
		ID: m["taskId"],
	}
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
	resp := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/db/import", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID), true, settings)
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
	resp := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/console", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID), true, settings)
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
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/services/%s/console/token/%s", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID, jobID), true, settings)
	var credentials models.ConsoleCredentials
	json.Unmarshal(resp, &credentials)
	return &credentials
}

// DestroyConsole properly shuts down a console service
func DestroyConsole(jobID string, serviceID string, settings *models.Settings) {
	httpclient.Delete(fmt.Sprintf("%s%s/environments/%s/services/%s/console/%s", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID, jobID), true, settings)
}

// RetrieveServiceMetrics fetches metrics for a single service for a specified
// number of minutes.
func RetrieveServiceMetrics(mins int, settings *models.Settings) *models.Metrics {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/metrics/%s?mins=%d", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, settings.ServiceID, mins), true, settings)
	var metrics models.Metrics
	json.Unmarshal(resp, &metrics)
	return &metrics
}

// ListServiceFiles retrieves a list of all downloadable service files for the
// specified code service.
func ListServiceFiles(serviceID string, settings *models.Settings) *[]models.ServiceFile {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/services/%s/files", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID), true, settings)
	var files []models.ServiceFile
	json.Unmarshal(resp, &files)
	return &files
}

// RetrieveServiceFile retrieves a service file by its ID.
func RetrieveServiceFile(serviceID string, fileID int64, settings *models.Settings) *models.ServiceFile {
	resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/services/%s/files/%d", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID, fileID), true, settings)
	var file models.ServiceFile
	json.Unmarshal(resp, &file)
	return &file
}
