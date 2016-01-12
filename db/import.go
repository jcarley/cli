package db

import (
	"crypto/aes"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
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
func Import(databaseLabel string, filePath string, mongoCollection string, mongoDatabase string, wipeFirst bool, settings *models.Settings) {
	helpers.SignIn(settings)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("A file does not exist at path '%s'\n", filePath)
		os.Exit(1)
	}
	service := helpers.RetrieveServiceByLabel(databaseLabel, settings)
	if service == nil {
		fmt.Printf("Could not find a service with the label \"%s\"\n", databaseLabel)
		os.Exit(1)
	}
	// backup before we do an import. so we can be safe
	fmt.Printf("Backing up \"%s\" before performing the import\n", databaseLabel)
	task := helpers.CreateBackup(service.ID, settings)
	fmt.Printf("Backup started (task ID = %s)\n", task.ID)
	fmt.Print("Polling until backup finishes.")
	ch := make(chan string, 1)
	go helpers.PollTaskStatus(task.ID, ch, settings)
	status := <-ch
	task.Status = status
	fmt.Printf("\nEnded in status '%s'\n", task.Status)
	helpers.DumpLogs(service, task, "backup", settings)
	if task.Status != "finished" {
		panic(fmt.Errorf("Backup finished in an invalid status '%s'\n", status))
	}
	// end backup section
	env := helpers.RetrieveEnvironment("spec", settings)
	pod := helpers.RetrievePodMetadata(env.PodID, settings)
	fmt.Printf("Importing '%s' into %s (ID = %s)\n", filePath, databaseLabel, service.ID)
	key := make([]byte, 32)
	iv := make([]byte, aes.BlockSize)
	rand.Read(key)
	rand.Read(iv)
	fmt.Println("Encrypting...")
	encrFilePath := helpers.EncryptFile(filePath, key, iv, pod.ImportRequiresLength)
	defer os.Remove(encrFilePath)
	options := map[string]string{}
	if mongoCollection != "" {
		options["mongoCollection"] = mongoCollection
	}
	if mongoDatabase != "" {
		options["mongoDatabase"] = mongoDatabase
	}
	fmt.Println("Uploading...")
	tempURL := helpers.RetrieveTempUploadURL(service.ID, settings)

	task = helpers.InitiateImport(tempURL.URL, encrFilePath, string(helpers.Base64Encode(helpers.Hex(key))), string(helpers.Base64Encode(helpers.Hex(iv))), options, wipeFirst, service.ID, settings)
	fmt.Printf("Processing import... (task ID = %s)\n", task.ID)

	ch = make(chan string, 1)
	go helpers.PollTaskStatus(task.ID, ch, settings)
	status = <-ch
	task.Status = status
	fmt.Printf("\nImport complete (end status = '%s')\n", task.Status)
	helpers.DumpLogs(service, task, "restore", settings)
	if task.Status != "finished" {
		os.Exit(1)
	}
}
