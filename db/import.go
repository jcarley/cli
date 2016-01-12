package db

import (
	"crypto/aes"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/catalyzeio/cli/helpers"
)

// Import imports data into a database service. The import is accomplished
// by encrypting the file locally, requesting a location that it can be uploaded
// to, then uploads the file. Once uploaded an automated service processes the
// file and acts according to the given parameters.
//
// The type of file that should be imported depends on the database. For
// PostgreSQL and MySQL, this should be a single `.sql` file. For Mongo, this
// should be a single tar'ed, gzipped archive (`.tar.gz`) of the database dump
// that you want to import.
func (d *SDb) Import() error {
	if _, err := os.Stat(d.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("A file does not exist at path '%s'\n", d.FilePath)
	}
	service := helpers.RetrieveServiceByLabel(d.DatabaseName, d.Settings)
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"\n", d.DatabaseName)
	}
	// backup before we do an import. so we can be safe
	fmt.Printf("Backing up \"%s\" before performing the import\n", d.DatabaseName)
	task := helpers.CreateBackup(service.ID, d.Settings)
	fmt.Printf("Backup started (task ID = %s)\n", task.ID)
	fmt.Print("Polling until backup finishes.")
	ch := make(chan string, 1)
	go helpers.PollTaskStatus(task.ID, ch, d.Settings)
	status := <-ch
	task.Status = status
	fmt.Printf("\nEnded in status '%s'\n", task.Status)
	helpers.DumpLogs(service, task, "backup", d.Settings)
	if task.Status != "finished" {
		return fmt.Errorf("Backup finished in an invalid status '%s'\n", status)
	}
	// end backup section
	env := helpers.RetrieveEnvironment("spec", d.Settings)
	pod := helpers.RetrievePodMetadata(env.PodID, d.Settings)
	fmt.Printf("Importing '%s' into %s (ID = %s)\n", d.FilePath, d.DatabaseName, service.ID)
	key := make([]byte, 32)
	iv := make([]byte, aes.BlockSize)
	rand.Read(key)
	rand.Read(iv)
	fmt.Println("Encrypting...")
	encrFilePath := helpers.EncryptFile(d.FilePath, key, iv, pod.ImportRequiresLength)
	defer os.Remove(encrFilePath)
	options := map[string]string{}
	if d.MongoCollection != "" {
		options["mongoCollection"] = d.MongoCollection
	}
	if d.MongoDatabase != "" {
		options["mongoDatabase"] = d.MongoDatabase
	}
	fmt.Println("Uploading...")
	tempURL := helpers.RetrieveTempUploadURL(service.ID, d.Settings)

	wipeFirst := false
	task = helpers.InitiateImport(tempURL.URL, encrFilePath, string(helpers.Base64Encode(helpers.Hex(key))), string(helpers.Base64Encode(helpers.Hex(iv))), options, wipeFirst, service.ID, d.Settings)
	fmt.Printf("Processing import... (task ID = %s)\n", task.ID)

	ch = make(chan string, 1)
	go helpers.PollTaskStatus(task.ID, ch, d.Settings)
	status = <-ch
	task.Status = status
	fmt.Printf("\nImport complete (end status = '%s')\n", task.Status)
	helpers.DumpLogs(service, task, "restore", d.Settings)
	if task.Status != "finished" {
		return fmt.Errorf("Finished with invalid status %s\n", task.Status)
	}
	return nil
}
