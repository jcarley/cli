package db

import (
	"crypto/aes"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/services"
)

func CmdImport(databaseName, filePath, mongoCollection, mongoDatabase string, id IDb, is services.IServices) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("A file does not exist at path '%s'\n", filePath)
	}
	service, err := is.RetrieveByLabel(databaseName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"\n", databaseName)
	}
	fmt.Printf("Backing up \"%s\" before performing the import\n", databaseName)
	err = id.Backup(false, service)
	if err != nil {
		return err
	}
	fmt.Printf("Importing '%s' into %s (ID = %s)\n", filePath, databaseName, service.ID)
	return id.Import(filePath, mongoCollection, mongoDatabase, service)
}

// Import imports data into a database service. The import is accomplished
// by encrypting the file locally, requesting a location that it can be uploaded
// to, then uploads the file. Once uploaded an automated service processes the
// file and acts according to the given parameters.
//
// The type of file that should be imported depends on the database. For
// PostgreSQL and MySQL, this should be a single `.sql` file. For Mongo, this
// should be a single tar'ed, gzipped archive (`.tar.gz`) of the database dump
// that you want to import.
func (d *SDb) Import(filePath, mongoCollection, mongoDatabase string, service *models.Service) error {
	key := make([]byte, 32)
	iv := make([]byte, aes.BlockSize)
	rand.Read(key)
	rand.Read(iv)
	fmt.Println("Encrypting...")
	encrFilePath := helpers.EncryptFile(filePath, key, iv, true)
	defer os.Remove(encrFilePath)
	options := map[string]string{}
	if mongoCollection != "" {
		options["mongoCollection"] = mongoCollection
	}
	if mongoDatabase != "" {
		options["mongoDatabase"] = mongoDatabase
	}
	fmt.Println("Uploading...")
	tempURL := helpers.RetrieveTempUploadURL(service.ID, d.Settings)

	wipeFirst := false
	task := helpers.InitiateImport(tempURL.URL, encrFilePath, string(helpers.Base64Encode(helpers.Hex(key))), string(helpers.Base64Encode(helpers.Hex(iv))), options, wipeFirst, service.ID, d.Settings)
	fmt.Printf("Processing import... (task ID = %s)\n", task.ID)

	ch := make(chan string, 1)
	go helpers.PollTaskStatus(task.ID, ch, d.Settings)
	status := <-ch
	task.Status = status
	fmt.Printf("\nImport complete (end status = '%s')\n", task.Status)
	helpers.DumpLogs(service, task, "restore", d.Settings)
	if task.Status != "finished" {
		return fmt.Errorf("Finished with invalid status %s\n", task.Status)
	}
	return nil
}
