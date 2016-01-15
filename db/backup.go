package db

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/services"
)

func CmdBackup(databaseName string, skipPoll bool, id IDb, is services.IServices) error {
	service, err := is.RetrieveByLabel(databaseName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"\n", databaseName)
	}
	return id.Backup(skipPoll, service)
}

// Backup creates a new backup
func (d *SDb) Backup(skipPoll bool, service *models.Service) error {
	task := helpers.CreateBackup(service.ID, d.Settings)
	fmt.Printf("Backup started (task ID = %s)\n", task.ID)
	if !skipPoll {
		fmt.Print("Polling until backup finishes.")
		ch := make(chan string, 1)
		go helpers.PollTaskStatus(task.ID, ch, d.Settings)
		status := <-ch
		task.Status = status
		fmt.Printf("\nEnded in status '%s'\n", task.Status)
		helpers.DumpLogs(service, task, "backup", d.Settings)
		if task.Status != "finished" {
			return fmt.Errorf("Task finished with invalid status %s\n", task.Status)
		}
	}
	return nil
}
