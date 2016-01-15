package db

// RestoreBackup a database service from an existing backup
/*func RestoreBackup(serviceLabel string, backupID string, skipPoll bool, settings *models.Settings) {
	helpers.SignIn(settings)
	service := helpers.RetrieveServiceByLabel(serviceLabel, settings)
	if service == nil {
		fmt.Printf("Could not find a service with the label \"%s\"\n", serviceLabel)
		os.Exit(1)
	}
	task := helpers.RestoreBackup(service.ID, backupID, settings)
	fmt.Printf("Restoring (task ID = %s)\n", task.ID)
	if !skipPoll {
		fmt.Print("Polling until restore finishes.")
		ch := make(chan string, 1)
		go helpers.PollTaskStatus(task.ID, ch, settings)
		status := <-ch
		task.Status = status
		fmt.Printf("\nEnded in status '%s'\n", task.Status)
		helpers.DumpLogs(service, task, "restore", settings)
		if task.Status != "finished" {
			os.Exit(1)
		}
	}
}*/
