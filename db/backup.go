package db

import (
	"fmt"
	"os"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

// CreateBackup a new backup
func CreateBackup(serviceLabel string, skipPoll bool, settings *models.Settings) {
	helpers.SignIn(settings)
	service := helpers.RetrieveServiceByLabel(serviceLabel, settings)
	if service == nil {
		fmt.Printf("Could not find a service with the label \"%s\"\n", serviceLabel)
		os.Exit(1)
	}
	task := helpers.CreateBackup(service.ID, settings)
	fmt.Printf("Backup started (task ID = %s)\n", task.ID)
	if !skipPoll {
		fmt.Print("Polling until backup finishes.")
		ch := make(chan string, 1)
		go helpers.PollTaskStatus(task.ID, ch, settings)
		status := <-ch
		task.Status = status
		fmt.Printf("\nEnded in status '%s'\n", task.Status)
		helpers.DumpLogs(service, task, "backup", settings)
		if task.Status != "finished" {
			os.Exit(1)
		}
	}
}
