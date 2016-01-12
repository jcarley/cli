package db

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/catalyzeio/cli/helpers"
)

// Download an existing backup to the local machine. The backup is encrypted
// throughout the entire journey and then decrypted once it is stored locally.
func (d *SDb) Download() error {
	err := d.Prompts.PHI()
	if err != nil {
		return err
	}
	if !d.Force {
		if _, err := os.Stat(d.FilePath); err == nil {
			return fmt.Errorf("File already exists at path '%s'. Specify `--force` to overwrite\n", d.FilePath)
		}
	} else {
		os.Remove(d.FilePath)
	}
	service := helpers.RetrieveServiceByLabel(d.DatabaseName, d.Settings)
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"\n", d.DatabaseName)
	}
	job := helpers.RetrieveJob(d.BackupID, service.ID, d.Settings)
	if job.Type != "backup" || job.Status != "finished" {
		fmt.Println("Only 'finished' 'backup' jobs may be downloaded")
	}
	fmt.Printf("Downloading backup %s\n", d.BackupID)
	tempURL := helpers.RetrieveTempURL(d.BackupID, service.ID, d.Settings)
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
	helpers.DecryptFile(tmpFile.Name(), job.Backup.Key, job.Backup.IV, d.FilePath)
	fmt.Printf("%s backup downloaded successfully to %s\n", d.DatabaseName, d.FilePath)
	return nil
}
