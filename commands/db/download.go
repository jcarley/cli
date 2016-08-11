package db

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
)

func CmdDownload(databaseName, backupID, filePath string, force bool, id IDb, ip prompts.IPrompts, is services.IServices) error {
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
	err = id.Download(backupID, filePath, service)
	if err != nil {
		return err
	}
	logrus.Printf("%s backup downloaded successfully to %s", databaseName, filePath)
	logrus.Printf("You can also view logs for this backup with the \"catalyze db logs %s %s\" command", databaseName, backupID)
	return nil
}

// Download an existing backup to the local machine. The backup is encrypted
// throughout the entire journey and then decrypted once it is stored locally.
func (d *SDb) Download(backupID, filePath string, service *models.Service) error {
	job, err := d.Jobs.Retrieve(backupID, service.ID, false)
	if err != nil {
		return err
	}
	if job.Type != "backup" || (job.Status != "finished" && job.Status != "disappeared") {
		return errors.New("Only 'finished' 'backup' jobs may be downloaded")
	}
	logrus.Printf("Downloading backup %s", backupID)
	tmpAuth, err := d.TempDownloadAuth(backupID, service)
	if err != nil {
		return err
	}
	u, err := url.Parse(tmpAuth.URL)
	if err != nil {
		return err
	}
	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1"), Credentials: credentials.NewStaticCredentials(tmpAuth.AccessKeyID, tmpAuth.SecretAccessKey, tmpAuth.SessionToken)}))
	decryptWriter, err := d.Crypto.NewDecryptFileWriterAt(filePath, job.Backup.Key, job.Backup.IV)
	if err != nil {
		return err
	}
	_, err = downloader.Download(decryptWriter, &s3.GetObjectInput{
		Bucket: aws.String(strings.Split(u.Host, ".")[0]),
		Key:    aws.String(strings.TrimLeft(u.Path, "/")),
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *SDb) TempDownloadAuth(jobID string, service *models.Service) (*models.TempAuth, error) {
	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod, d.Settings.UsersID)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/temp-backup-auth/%s", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, service.ID, jobID), headers)
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
