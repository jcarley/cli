package db

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

// Export dumps all data from a database service and downloads the encrypted
// data to the local machine. The export is accomplished by first creating a
// backup. Once finished, the CLI asks where the file can be downloaded from.
// The file is downloaded, decrypted, and saved locally.
func Export(databaseLabel string, filePath string, force bool, settings *models.Settings) {
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
	service := helpers.RetrieveServiceByLabel(databaseLabel, settings)
	if service == nil {
		fmt.Printf("Could not find a service with the label \"%s\"\n", databaseLabel)
		os.Exit(1)
	}
	task := helpers.CreateBackup(service.ID, settings)
	fmt.Printf("Export started (task ID = %s)\n", task.ID)
	fmt.Print("Polling until export finishes.")
	ch := make(chan string, 1)
	go helpers.PollTaskStatus(task.ID, ch, settings)
	status := <-ch
	task.Status = status
	if task.Status != "finished" {
		fmt.Printf("\nExport finished with illegal status \"%s\", aborting.\n", task.Status)
		helpers.DumpLogs(service, task, "backup", settings)
		os.Exit(1)
	}
	fmt.Printf("\nEnded in status '%s'\n", task.Status)
	job := helpers.RetrieveJobFromTaskID(task.ID, settings)
	fmt.Printf("Downloading export %s\n", job.ID)
	tempURL := helpers.RetrieveTempURL(job.ID, service.ID, settings)
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
	fmt.Println("Decrypting...")
	tmpFile.Close()
	helpers.DecryptFile(tmpFile.Name(), job.Backup.Key, job.Backup.IV, filePath)
	fmt.Printf("%s exported successfully to %s\n", databaseLabel, filePath)
	helpers.DumpLogs(service, task, "backup", settings)
}
