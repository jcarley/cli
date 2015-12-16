package helpers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/catalyzeio/catalyze/config"
	"github.com/catalyzeio/catalyze/httpclient"
	"github.com/catalyzeio/catalyze/models"
)

// PollTaskStatus polls a task until its status is not `scheduled`, `queued`,
// `started`, or `running`. The tasks status is then sent back through the
// chan.
func PollTaskStatus(taskID string, ch chan string, settings *models.Settings) {
	var task models.Task
poll:
	for {
		resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/tasks/%s", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, taskID), true, settings)
		json.Unmarshal(resp, &task)
		switch task.Status {
		case "scheduled", "queued", "started", "running":
			fmt.Print(".")
			time.Sleep(2 * time.Second)
		case "finished":
			break poll
		default:
			fmt.Printf("Error - ended in status '%s'.\n", task.Status)
			break poll
		}
	}
	ch <- task.Status
}

// PollConsoleJob polls a console job until it gets a jobId back
func PollConsoleJob(taskID string, serviceID string, ch chan string, settings *models.Settings) {
	job := make(map[string]string)
poll:
	for {
		resp := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/services/%s/console/status/%s", settings.PaasHost, config.PaasHostVersion, settings.EnvironmentID, serviceID, taskID), true, settings)
		json.Unmarshal(resp, &job)
		if jobID, ok := job["jobId"]; ok && jobID != "" {
			break poll
		} else {
			fmt.Print(".")
			time.Sleep(2 * time.Second)
		}
	}
	ch <- job["jobId"]
}
