package tasks

import (
	"fmt"
	"time"

	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
)

// PollForStatus polls a task until its status is not `scheduled`, `queued`,
// `started`, or `running`. The tasks status is then sent back through the
// chan.
func (t *STasks) PollForStatus(pollTask *models.Task) (string, error) {
	var task models.Task
	headers := httpclient.GetHeaders(t.Settings.APIKey, t.Settings.SessionToken, t.Settings.Version, t.Settings.Pod)
poll:
	for {
		resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/tasks/%s", t.Settings.PaasHost, t.Settings.PaasHostVersion, t.Settings.EnvironmentID, pollTask.ID), headers)
		if err != nil {
			return "", nil
		}
		err = httpclient.ConvertResp(resp, statusCode, &task)
		if err != nil {
			return "", nil
		}
		switch task.Status {
		case "scheduled", "queued", "started", "running":
			fmt.Print(".")
			time.Sleep(2 * time.Second)
		case "finished":
			break poll
		default:
			return "", fmt.Errorf("Error - ended in status '%s'.\n", task.Status)
		}
	}
	return task.Status, nil
}

// PollForConsole polls a console job until it gets a jobId back
func (t *STasks) PollForConsole(task *models.Task, service *models.Service) (string, error) {
	job := make(map[string]string)
	headers := httpclient.GetHeaders(t.Settings.APIKey, t.Settings.SessionToken, t.Settings.Version, t.Settings.Pod)
	for {
		resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/console/status/%s", t.Settings.PaasHost, t.Settings.PaasHostVersion, t.Settings.EnvironmentID, service.ID, task.ID), headers)
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
			fmt.Print(".")
			time.Sleep(2 * time.Second)
		}
	}
	return job["jobId"], nil
}
