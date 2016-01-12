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
