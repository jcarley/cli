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

func CmdDownload(databaseName, backupID, filePath string, force bool, id IDb, ip prompts.IPrompts, is services.IServices) error {
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
	err = id.Download(backupID, filePath, service)
	if err != nil {
		return err
	}
	fmt.Printf("%s backup downloaded successfully to %s\n", databaseName, filePath)
	return nil
}

// Download an existing backup to the local machine. The backup is encrypted
// throughout the entire journey and then decrypted once it is stored locally.
func (d *SDb) Download(backupID, filePath string, service *models.Service) error {
	job := helpers.RetrieveJob(backupID, service.ID, d.Settings)
	if job.Type != "backup" || job.Status != "finished" {
		fmt.Println("Only 'finished' 'backup' jobs may be downloaded")
	}
	fmt.Printf("Downloading backup %s\n", backupID)
	tempURL := helpers.RetrieveTempURL(backupID, service.ID, d.Settings)
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
	tmpFile.Close()
	fmt.Println("Decrypting...")
	helpers.DecryptFile(tmpFile.Name(), job.Backup.Key, job.Backup.IV, filePath)

	return nil
}
