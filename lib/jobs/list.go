package jobs

import (
	"fmt"

	"github.com/catalyzeio/cli/models"
)

func (j *SJobs) List(svcID string, page, pageSize int) (*[]models.Job, error) {
	headers := j.Settings.HTTPManager.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod, j.Settings.UsersID)
	resp, statusCode, err := j.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs?pageNumber=%d&pageSize=%d", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, svcID, page, pageSize), headers)
	if err != nil {
		return nil, err
	}
	var jobs []models.Job
	err = j.Settings.HTTPManager.ConvertResp(resp, statusCode, &jobs)
	if err != nil {
		return nil, err
	}
	return &jobs, nil
}
