package jobs

import (
	"fmt"

	"github.com/catalyzeio/cli/models"
)

func (j *SJobs) Retrieve(jobID, svcID string, includeSpec bool) (*models.Job, error) {
	headers := j.Settings.HTTPManager.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod, j.Settings.UsersID)
	resp, statusCode, err := j.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs/%s?spec=true", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, svcID, jobID), headers)
	if err != nil {
		return nil, err
	}
	var job models.Job
	err = j.Settings.HTTPManager.ConvertResp(resp, statusCode, &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (j *SJobs) RetrieveByStatus(svcID, status string) (*[]models.Job, error) {
	headers := j.Settings.HTTPManager.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod, j.Settings.UsersID)
	resp, statusCode, err := j.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs?status=%s", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, svcID, status), headers)
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

func (j *SJobs) RetrieveByType(svcID, jobType string, page, pageSize int) (*[]models.Job, error) {
	headers := j.Settings.HTTPManager.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod, j.Settings.UsersID)
	resp, statusCode, err := j.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs?type=%s&pageNumber=%d&pageSize=%d", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, svcID, jobType, page, pageSize), headers)
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

func (j *SJobs) RetrieveByTarget(svcID, target string, page, pageSize int) (*[]models.Job, error) {
	headers := j.Settings.HTTPManager.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod, j.Settings.UsersID)
	resp, statusCode, err := j.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs?type=worker&target=%s&pageNumber=%d&pageSize=%d", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, svcID, target, page, pageSize), headers)
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
