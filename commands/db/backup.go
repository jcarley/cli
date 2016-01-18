package db

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/lib/tasks"
	"github.com/catalyzeio/cli/models"
)

func CmdBackup(databaseName string, skipPoll bool, id IDb, is services.IServices, it tasks.ITasks) error {
	service, err := is.RetrieveByLabel(databaseName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"", databaseName)
	}
	task, err := id.Backup(service)
	if err != nil {
		return err
	}
	logrus.Printf("Backup started (task ID = %s)", task.ID)
	if !skipPoll {
		logrus.Print("Polling until backup finishes.")
		status, err := it.PollForStatus(task)
		if err != nil {
			return err
		}
		task.Status = status
		logrus.Printf("\nEnded in status '%s'", task.Status)
		err = id.DumpLogs("backup", task, service)
		if err != nil {
			return err
		}
		if task.Status != "finished" {
			return fmt.Errorf("Task finished with invalid status %s", task.Status)
		}
	}
	return nil
}

// Backup creates a new backup
func (d *SDb) Backup(service *models.Service) (*models.Task, error) {
	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod)
	resp, statusCode, err := httpclient.Post(nil, fmt.Sprintf("%s%s/services/%s/brrgc/backup", d.Settings.PaasHost, d.Settings.PaasHostVersion, service.ID), headers)
	if err != nil {
		return nil, err
	}
	var m map[string]string
	err = httpclient.ConvertResp(resp, statusCode, &m)
	if err != nil {
		return nil, err
	}
	return &models.Task{
		ID: m["task"],
	}, nil
}
