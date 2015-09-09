package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
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

// DownloadBackup an existing backup to the local machine. The backup is encrypted
// throughout the entire journey and then decrypted once it is stored locally.
func DownloadBackup(serviceLabel string, backupID string, filePath string, force bool, settings *models.Settings) {
	helpers.PHIPrompt()
	helpers.SignIn(settings)
	if !force {
		if _, err := os.Stat(filePath); err == nil {
			fmt.Printf("File already exists at path '%s'. Specify `--force` to overwrite\n", filePath)
			os.Exit(1)
		}
	} else {
		os.Remove(filePath)
	}
	service := helpers.RetrieveServiceByLabel(serviceLabel, settings)
	if service == nil {
		fmt.Printf("Could not find a service with the label \"%s\"\n", serviceLabel)
		os.Exit(1)
	}
	job := helpers.RetrieveJob(backupID, service.ID, settings)
	if job.Type != "backup" || job.Status != "finished" {
		fmt.Println("Only 'finished' 'backup' jobs may be downloaded")
	}
	fmt.Printf("Downloading backup %s\n", backupID)
	tempURL := helpers.RetrieveTempURL(backupID, service.ID, settings)
	dir, dirErr := ioutil.TempDir("", "")
	if dirErr != nil {
		fmt.Println(dirErr.Error())
		os.Exit(1)
	}
	defer os.Remove(dir)
	tmpFile, tmpFileErr := ioutil.TempFile(dir, "")
	if tmpFileErr != nil {
		fmt.Println(tmpFileErr.Error())
		os.Exit(1)
	}
	resp, respErr := http.Get(tempURL.URL)
	if respErr != nil {
		fmt.Println(respErr.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	io.Copy(tmpFile, resp.Body)
	tmpFile.Close()
	fmt.Println("Decrypting...")
	helpers.DecryptFile(tmpFile.Name(), job.Backup.Key, job.Backup.IV, filePath)
	fmt.Printf("%s backup downloaded successfully to %s\n", serviceLabel, filePath)
}

// SortedJobs is a wrapper for Jobs array in order to sort them by CreatedAt
// for the ListBackups command
type SortedJobs []models.Job

func (jobs SortedJobs) Len() int {
	return len(jobs)
}

func (jobs SortedJobs) Swap(i, j int) {
	jobs[i], jobs[j] = jobs[j], jobs[i]
}

func (jobs SortedJobs) Less(i, j int) bool {
	return jobs[i].CreatedAt < jobs[j].CreatedAt
}

// ListBackups lists the created backups for the service sorted from oldest to newest
func ListBackups(serviceLabel string, page int, pageSize int, settings *models.Settings) {
	helpers.SignIn(settings)
	service := helpers.RetrieveServiceByLabel(serviceLabel, settings)
	if service == nil {
		fmt.Printf("Could not find a service with the label \"%s\"\n", serviceLabel)
		os.Exit(1)
	}
	jobs := helpers.ListBackups(service.ID, page, pageSize, settings)
	sort.Sort(SortedJobs(*jobs))
	for _, job := range *jobs {
		fmt.Printf("%s %s (status = %s)\n", job.ID, job.CreatedAt, job.Status)
	}
	if len(*jobs) == pageSize && page == 1 {
		fmt.Println("(for older backups, try with --page 2 or adjust --page-size)")
	}
	if len(*jobs) == 0 && page == 1 {
		fmt.Println("No backups created yet for this service.")
	}
}

// RestoreBackup a database service from an existing backup
func RestoreBackup(serviceLabel string, backupID string, skipPoll bool, settings *models.Settings) {
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
}
