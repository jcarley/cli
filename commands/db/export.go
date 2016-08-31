package db

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/tj/go-spin"
)

func CmdExport(databaseName, filePath string, force bool, id IDb, ip prompts.IPrompts, is services.IServices, ij jobs.IJobs) error {
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
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services\" command.", databaseName)
	}
	job, err := id.Backup(service)
	if err != nil {
		return err
	}
	logrus.Printf("Export started (job ID = %s)", job.ID)
	// all because logrus treats print, println, and printf the same
	logrus.StandardLogger().Out.Write([]byte("Polling until export finishes."))
	status, err := ij.PollTillFinished(job.ID, service.ID)
	if err != nil {
		return err
	}
	job.Status = status
	logrus.Printf("\nEnded in status '%s'", job.Status)
	if job.Status != "finished" {
		id.DumpLogs("backup", job, service)
		return fmt.Errorf("Job finished with invalid status %s", job.Status)
	}

	err = id.Export(filePath, job, service)
	if err != nil {
		return err
	}
	err = id.DumpLogs("backup", job, service)
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
func (d *SDb) Export(filePath string, job *models.Job, service *models.Service) error {
	spinner := spin.New()
	done := make(chan struct{}, 1)
	defer func() { done <- struct{}{} }()
	go func() {
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				fmt.Printf("\r\033[mDownloading export %s. This may take awhile %s\033[m ", job.ID, spinner.Next())
			case <-done:
				return
			}
		}
	}()
	tempURL, err := d.TempDownloadURL(job.ID, service)
	if err != nil {
		done <- struct{}{}
		return err
	}
	resp, err := http.Get(tempURL.URL)
	if err != nil {
		done <- struct{}{}
		return err
	}
	defer resp.Body.Close()
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		done <- struct{}{}
		return err
	}
	dfw, err := d.Crypto.NewDecryptWriteCloser(file, job.Backup.Key, job.Backup.IV)
	if err != nil {
		done <- struct{}{}
		return err
	}
	_, err = io.Copy(dfw, resp.Body)
	if err != nil {
		return err
	}
	return dfw.Close()
}
