package jobs

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func (j *SJobs) PollForStatus(jobID, svcID string) (string, error) {
	var job models.Job
	failedAttempts := 0
poll:
	for {
		failed := false
		headers := httpclient.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod)
		resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs/%s", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, svcID, jobID), headers)
		if err != nil {
			failed = true
		}
		err = httpclient.ConvertResp(resp, statusCode, &job)
		if err != nil {
			failed = true
		}
		if failed {
			failedAttempts++
		}
		switch job.Status {
		case "scheduled", "queued", "started", "running":
			if failedAttempts >= 3 {
				return "", fmt.Errorf("Error - ended in status '%s'.", job.Status)
			}
			logrus.Print(".")
			time.Sleep(config.JobPollTime * time.Second)
		case "finished":
			break poll
		default:
			return "", fmt.Errorf("Error - ended in status '%s'.", job.Status)
		}
	}
	return job.Status, nil
}
