package db

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
)

// Backup creates a new backup
func (d *SDb) Backup() error {
	helpers.SignIn(d.Settings)
	service := helpers.RetrieveServiceByLabel(d.DatabaseName, d.Settings)
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"\n", d.DatabaseName)
	}
	task := helpers.CreateBackup(service.ID, d.Settings)
	fmt.Printf("Backup started (task ID = %s)\n", task.ID)
	if !d.SkipPoll {
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
