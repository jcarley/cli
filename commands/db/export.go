package db

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/lib/tasks"
	"github.com/catalyzeio/cli/models"
)

func CmdExport(databaseName, filePath string, force bool, id IDb, ip prompts.IPrompts, is services.IServices, it tasks.ITasks) error {
	err := ip.PHI()
	if err != nil {
		return err
	}
	if !force {
		if _, err := os.Stat(filePath); err == nil {
			return fmt.Errorf("File already exists at path '%s'. Specify `--force` to overwrite", filePath)
		}
	} else {
		os.Remove(filePath)
	}
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
	logrus.Printf("Export started (task ID = %s)", task.ID)
	logrus.Print("Polling until export finishes.")
	status, err := it.PollForStatus(task)
	if err != nil {
		return err
	}
	task.Status = status
	logrus.Printf("\nEnded in status '%s'", task.Status)
	if task.Status != "finished" {
		id.DumpLogs("backup", task, service)
		return fmt.Errorf("Task finished with invalid status %s", task.Status)
	}

	err = id.Export(filePath, task, service)
	if err != nil {
		return err
	}
	err = id.DumpLogs("backup", task, service)
	if err != nil {
		return err
	}
	logrus.Printf("%s exported successfully to %s", service.Name, filePath)
	return nil
}

// Export dumps all data from a database service and downloads the encrypted
// data to the local machine. The export is accomplished by first creating a
// backup. Once finished, the CLI asks where the file can be downloaded from.
// The file is downloaded, decrypted, and saved locally.
func (d *SDb) Export(filePath string, task *models.Task, service *models.Service) error {
	job, err := d.Jobs.RetrieveFromTaskID(task.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Downloading export %s", job.ID)
	tempURL, err := d.TempDownloadURL(job.ID, service)
	if err != nil {
		return err
	}
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.Remove(dir)
	tmpFile, err := ioutil.TempFile(dir, "")
	if err != nil {
		return err
	}
	resp, err := http.Get(tempURL.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(tmpFile, resp.Body)
	logrus.Println("Decrypting...")
	tmpFile.Close()
	err = d.Crypto.DecryptFile(tmpFile.Name(), job.Backup.Key, job.Backup.IV, filePath)
	if err != nil {
		return err
	}
	return nil
}
