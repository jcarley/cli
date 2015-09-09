package helpers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/catalyzeio/catalyze/models"
)

// DumpLogs dumps logs from a Backup/Restore/Import/Export job to the console
func DumpLogs(service *models.Service, task *models.Task, taskType string, settings *models.Settings) {
	fmt.Printf("Retrieving %s logs for task %s...\n", service.Label, task.ID)
	job := RetrieveJobFromTaskID(task.ID, settings)
	tempURL := RetrieveTempLogsURL(job.ID, taskType, service.ID, settings)
	dir, dirErr := ioutil.TempDir("", "")
	if dirErr != nil {
		fmt.Println(dirErr.Error())
		os.Exit(1)
	}

	encrFile, encrFileErr := ioutil.TempFile(dir, "")
	if encrFileErr != nil {
		fmt.Println(encrFileErr.Error())
		os.Exit(1)
	}
	resp, respErr := http.Get(tempURL.URL)
	if respErr != nil {
		fmt.Println(respErr.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	io.Copy(encrFile, resp.Body)
	encrFile.Close()

	plainFile, plainFileErr := ioutil.TempFile(dir, "")
	if plainFileErr != nil {
		fmt.Println(plainFileErr.Error())
		os.Exit(1)
	}
	// do we have to close this before calling DecryptFile?
	// or can two processes have a file open simultaneously?
	plainFile.Close()

	if taskType == "backup" {
		DecryptFile(encrFile.Name(), job.Backup.Key, job.Backup.IV, plainFile.Name())
	} else if taskType == "restore" {
		DecryptFile(encrFile.Name(), job.Restore.Key, job.Restore.IV, plainFile.Name())
	}
	fmt.Printf("-------------------------- Begin %s logs --------------------------\n", service.Label)
	plainFile, _ = os.Open(plainFile.Name())
	io.Copy(os.Stdout, plainFile)
	plainFile.Close()
	fmt.Printf("--------------------------- End %s logs ---------------------------\n", service.Label)
	os.Remove(encrFile.Name())
	os.Remove(plainFile.Name())
	os.Remove(dir)
}
