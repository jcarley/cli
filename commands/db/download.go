package db

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/tj/go-spin"
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
	spinner := spin.New()
	done := make(chan struct{}, 1)
	defer func() { done <- struct{}{} }()
	go func() {
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				fmt.Printf("\r\033[mDownloading backup %s. This may take awhile %s\033[m ", backupID, spinner.Next())
			case <-done:
				return
			}
		}
	}()
	tempURL, err := d.TempDownloadURL(backupID, service)
	if err != nil {
		done <- struct{}{}
		return err
	}

	u, _ := url.Parse(tempURL.URL)
	svc := s3.New(session.New(&aws.Config{Region: aws.String("us-east-1"), Credentials: credentials.AnonymousCredentials}))
	req, resp := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(strings.Split(u.Host, ".")[0]),
		Key:    aws.String(strings.TrimLeft(u.Path, "/")),
	})
	req.HTTPRequest.URL.RawQuery = u.RawQuery
	err = req.Send()
	if err != nil {
		done <- struct{}{}
		return err
	}
	defer resp.Body.Close()
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	dfw, err := d.Crypto.NewDecryptWriteCloser(file, job.Backup.Key, job.Backup.IV)
	if err != nil {
		return err
	}
	_, err = io.Copy(dfw, resp.Body)
	if err != nil {
		return err
	}
	return dfw.Close()
}

func (d *SDb) TempDownloadURL(jobID string, service *models.Service) (*models.TempURL, error) {
	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod, d.Settings.UsersID)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/backup-url/%s", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, service.ID, jobID), headers)
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
