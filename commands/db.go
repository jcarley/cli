package commands

import (
	"crypto/aes"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
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
		os.Exit(1)
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
