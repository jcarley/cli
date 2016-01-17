package db

import (
	"encoding/json"
	"fmt"

	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/services"
	"github.com/catalyzeio/cli/tasks"
)

func CmdBackup(databaseName string, skipPoll bool, id IDb, is services.IServices, it tasks.ITasks) error {
	service, err := is.RetrieveByLabel(databaseName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"\n", databaseName)
	}
	task, err := id.Backup(service)
	if err != nil {
		return err
	}
	fmt.Printf("Backup started (task ID = %s)\n", task.ID)
	if !skipPoll {
		fmt.Print("Polling until backup finishes.")
		status, err := it.PollForStatus(task)
		if err != nil {
			return err
		}
		task.Status = status
		fmt.Printf("\nEnded in status '%s'\n", task.Status)
		err = id.DumpLogs("backup", task, service)
		if err != nil {
			return err
		}
		if task.Status != "finished" {
			return fmt.Errorf("Task finished with invalid status %s\n", task.Status)
		}
	}
	return nil
}

// Backup creates a new backup
func (d *SDb) Backup(service *models.Service) (*models.Task, error) {
	backup := map[string]string{
		"archiveType":    "cf",
		"encryptionType": "aes",
	}
	b, err := json.Marshal(backup)
	if err != nil {
		return nil, err
	}
	headers := httpclient.GetHeaders(d.Settings.APIKey, d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/backup", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, service.ID), headers)
	if err != nil {
		return nil, err
	}
	var m map[string]string
	err = httpclient.ConvertResp(resp, statusCode, &m)
	return &models.Task{
		ID: m["taskId"],
	}, nil
}
