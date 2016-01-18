package db

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"

	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/services"
	"github.com/catalyzeio/cli/tasks"
)

func CmdImport(databaseName, filePath, mongoCollection, mongoDatabase string, id IDb, is services.IServices, it tasks.ITasks) error {
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
	task, err := id.Backup(service)
	if err != nil {
		return err
	}
	fmt.Printf("Backup started (task ID = %s)\n", task.ID)

	fmt.Print("Polling until backup finishes.")
	status, err := it.PollForStatus(task)
	if err != nil {
		return err
	}
	task.Status = status
	fmt.Printf("\nEnded in status '%s'\n", task.Status)
	err = id.DumpLogs("backup", task, service)
	if err != nil {
		return err
	}
	if task.Status != "finished" {
		return fmt.Errorf("Task finished with invalid status %s\n", task.Status)
	}
	fmt.Printf("Importing '%s' into %s (ID = %s)\n", filePath, databaseName, service.ID)
	task, err = id.Import(filePath, mongoCollection, mongoDatabase, service)
	if err != nil {
		return err
	}
	fmt.Printf("Processing import... (task ID = %s)\n", task.ID)

	status, err = it.PollForStatus(task)
	if err != nil {
		return err
	}
	task.Status = status
	fmt.Printf("\nImport complete (end status = '%s')\n", task.Status)
	err = id.DumpLogs("restore", task, service)
	if err != nil {
		return err
	}
	if task.Status != "finished" {
		return fmt.Errorf("Finished with invalid status %s\n", task.Status)
	}
	return nil
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
func (d *SDb) Import(filePath, mongoCollection, mongoDatabase string, service *models.Service) (*models.Task, error) {
	key := make([]byte, 32)
	iv := make([]byte, aes.BlockSize)
	rand.Read(key)
	rand.Read(iv)
	fmt.Println("Encrypting...")
	encrFilePath, err := d.Crypto.EncryptFile(filePath, key, iv)
	if err != nil {
		return nil, err
	}
	defer os.Remove(encrFilePath)
	options := map[string]string{}
	if mongoCollection != "" {
		options["mongoCollection"] = mongoCollection
	}
	if mongoDatabase != "" {
		options["mongoDatabase"] = mongoDatabase
	}
	fmt.Println("Uploading...")
	tempURL, err := d.TempUploadURL(service)
	if err != nil {
		return nil, err
	}

	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod)
	resp, statusCode, err := httpclient.PutFile(encrFilePath, tempURL.URL, headers)
	if err != nil {
		return nil, err
	}
	err = httpclient.ConvertResp(resp, statusCode, nil)
	if err != nil {
		return nil, err
	}
	importParams := map[string]interface{}{}
	for key, value := range options {
		importParams[key] = value
	}
	importParams["filename"] = tempURL.URL
	importParams["encryptionKey"] = string(d.Crypto.Base64Encode(d.Crypto.Hex(key)))
	importParams["encryptionIV"] = string(d.Crypto.Base64Encode(d.Crypto.Hex(iv)))
	importParams["dropDatabase"] = false

	b, err := json.Marshal(importParams)
	if err != nil {
		return nil, err
	}
	resp, statusCode, err = httpclient.Post(b, fmt.Sprintf("%s%s/services/%s/import", d.Settings.PaasHost, d.Settings.PaasHostVersion, service.ID), headers)
	if err != nil {
		return nil, err
	}
	var m map[string]string
	err = httpclient.ConvertResp(resp, statusCode, &m)
	if err != nil {
		return nil, err
	}
	return &models.Task{
		ID: m["task"],
	}, nil
}

func (d *SDb) TempUploadURL(service *models.Service) (*models.TempURL, error) {
	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/services/%s/restore-url", d.Settings.PaasHost, d.Settings.PaasHostVersion, service.ID), headers)
	if err != nil {
		return nil, err
	}
	var tempURL models.TempURL
	err = httpclient.ConvertResp(resp, statusCode, &tempURL)
	if err != nil {
		return nil, err
	}
	return &tempURL, nil
}
