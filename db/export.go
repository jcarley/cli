package db

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/prompts"
	"github.com/catalyzeio/cli/services"
)

func CmdExport(databaseName, filePath string, force bool, id IDb, ip prompts.IPrompts, is services.IServices) error {
	err := ip.PHI()
	if err != nil {
		return err
	}
	if !force {
		if _, err := os.Stat(filePath); err == nil {
			return fmt.Errorf("File already exists at path '%s'. Specify `--force` to overwrite\n", filePath)
		}
	} else {
		os.Remove(filePath)
	}
	service, err := is.RetrieveByLabel(databaseName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"\n", databaseName)
	}
	return id.Export(filePath, service)
}

// Export dumps all data from a database service and downloads the encrypted
// data to the local machine. The export is accomplished by first creating a
// backup. Once finished, the CLI asks where the file can be downloaded from.
// The file is downloaded, decrypted, and saved locally.
func (d *SDb) Export(filePath string, service *models.Service) error {
	task := helpers.CreateBackup(service.ID, d.Settings)
	fmt.Printf("Export started (task ID = %s)\n", task.ID)
	fmt.Print("Polling until export finishes.")
	ch := make(chan string, 1)
	go helpers.PollTaskStatus(task.ID, ch, d.Settings)
	status := <-ch
	task.Status = status
	if task.Status != "finished" {
		helpers.DumpLogs(service, task, "backup", d.Settings)
		return fmt.Errorf("\nExport finished with illegal status \"%s\", aborting.\n", task.Status)
	}
	fmt.Printf("\nEnded in status '%s'\n", task.Status)
	job := helpers.RetrieveJobFromTaskID(task.ID, d.Settings)
	fmt.Printf("Downloading export %s\n", job.ID)
	tempURL := helpers.RetrieveTempURL(job.ID, service.ID, d.Settings)
	dir, dirErr := ioutil.TempDir("", "")
	if dirErr != nil {
		return dirErr
	}
	defer os.Remove(dir)
	tmpFile, tmpFileErr := ioutil.TempFile(dir, "")
	if tmpFileErr != nil {
		return tmpFileErr
	}
	resp, respErr := http.Get(tempURL.URL)
	if respErr != nil {
		return respErr
	}
	defer resp.Body.Close()
	io.Copy(tmpFile, resp.Body)
	fmt.Println("Decrypting...")
	tmpFile.Close()
	helpers.DecryptFile(tmpFile.Name(), job.Backup.Key, job.Backup.IV, filePath)
	fmt.Printf("%s exported successfully to %s\n", service.Name, filePath)
	helpers.DumpLogs(service, task, "backup", d.Settings)
	return nil
}
