package db

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/crypto"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/models"
)

func CmdImport(databaseName, filePath, mongoCollection, mongoDatabase string, id IDb, is services.IServices, ij jobs.IJobs) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("A file does not exist at path '%s'", filePath)
	}
	service, err := is.RetrieveByLabel(databaseName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services\" command.", databaseName)
	}
	logrus.Printf("Backing up \"%s\" before performing the import", databaseName)
	job, err := id.Backup(service)
	if err != nil {
		return err
	}
	logrus.Printf("Backup started (job ID = %s)", job.ID)

	// all because logrus treats print, println, and printf the same
	logrus.StandardLogger().Out.Write([]byte("Polling until backup finishes."))
	status, err := ij.PollTillFinished(job.ID, service.ID)
	if err != nil {
		return err
	}
	job.Status = status
	logrus.Printf("\nEnded in status '%s'", job.Status)
	err = id.DumpLogs("backup", job, service)
	if err != nil {
		return err
	}
	if job.Status != "finished" {
		return fmt.Errorf("Job finished with invalid status %s", job.Status)
	}
	logrus.Printf("Importing '%s' into %s (ID = %s)", filePath, databaseName, service.ID)
	job, err = id.Import(filePath, mongoCollection, mongoDatabase, service)
	if err != nil {
		return err
	}
	// all because logrus treats print, println, and printf the same
	logrus.StandardLogger().Out.Write([]byte(fmt.Sprintf("Processing import (job ID = %s).", job.ID)))

	status, err = ij.PollTillFinished(job.ID, service.ID)
	if err != nil {
		return err
	}
	job.Status = status
	logrus.Printf("\nImport complete (end status = '%s')", job.Status)
	err = id.DumpLogs("restore", job, service)
	if err != nil {
		return err
	}
	if job.Status != "finished" {
		return fmt.Errorf("Finished with invalid status %s", job.Status)
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
func (d *SDb) Import(filePath, mongoCollection, mongoDatabase string, service *models.Service) (*models.Job, error) {
	key := make([]byte, crypto.KeySize)
	iv := make([]byte, crypto.IVSize)
	rand.Read(key)
	rand.Read(iv)
	encryptFileReader, err := d.Crypto.NewEncryptFileReader(filePath, key, iv)
	defer encryptFileReader.Close()
	if err != nil {
		return nil, err
	}
	options := map[string]string{}
	if mongoCollection != "" {
		options["databaseCollection"] = mongoCollection
	}
	if mongoDatabase != "" {
		options["database"] = mongoDatabase
	}
	logrus.Println("Encrypting and Uploading...")
	tmpAuth, err := d.TempUploadAuth(service)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(tmpAuth.URL)
	if err != nil {
		return nil, err
	}
	fileName := strings.TrimLeft(u.Path, "/")
	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("us-east-1"), Credentials: credentials.NewStaticCredentials(tmpAuth.AccessKeyID, tmpAuth.SecretAccessKey, tmpAuth.SessionToken)}))
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(strings.Split(u.Host, ".")[0]),
		Key:                  aws.String(fileName),
		Body:                 encryptFileReader,
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return nil, err
	}
	importParams := map[string]interface{}{}
	for key, value := range options {
		importParams[key] = value
	}
	importParams["filename"] = fileName
	importParams["encryptionKey"] = string(d.Crypto.Hex(key, crypto.KeySize*2))
	importParams["encryptionIV"] = string(d.Crypto.Hex(iv, crypto.IVSize*2))
	importParams["dropDatabase"] = false

	b, err := json.Marshal(importParams)
	if err != nil {
		return nil, err
	}
	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod, d.Settings.UsersID)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/import", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, service.ID), headers)
	if err != nil {
		return nil, err
	}
	var job models.Job
	err = httpclient.ConvertResp(resp, statusCode, &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (d *SDb) TempUploadAuth(service *models.Service) (*models.TempAuth, error) {
	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod, d.Settings.UsersID)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/temp-restore-auth", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, service.ID), headers)
	if err != nil {
		return nil, err
	}
	var tempAuth models.TempAuth
	err = httpclient.ConvertResp(resp, statusCode, &tempAuth)
	if err != nil {
		return nil, err
	}
	return &tempAuth, nil
}
