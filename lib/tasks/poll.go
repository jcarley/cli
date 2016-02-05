package tasks

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

// PollForConsole polls a console job until it gets a jobId back
func (t *STasks) PollForConsole(task *models.Task, service *models.Service) (string, error) {
	job := make(map[string]string)
	headers := httpclient.GetHeaders(t.Settings.SessionToken, t.Settings.Version, t.Settings.Pod)
	for {
		resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/console/status/%s", t.Settings.PaasHost, t.Settings.PaasHostVersion, task.ID), headers)
		if err != nil {
			return "", err
		}
		err = httpclient.ConvertResp(resp, statusCode, &job)
		if err != nil {
			return "", err
		}
		if jobID, ok := job["jobId"]; ok && jobID != "" {
			break
		} else {
			logrus.Print(".")
			time.Sleep(config.JobPollTime * time.Second)
		}
	}
	return job["jobId"], nil
}
