package jobs

import (
	"fmt"

	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func (j *SJobs) Retrieve(jobID, svcID string, includeSpec bool) (*models.Job, error) {
	headers := httpclient.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod, j.Settings.UsersID)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs/%s?spec=true", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, svcID, jobID), headers)
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

func (j *SJobs) RetrieveByStatus(svcID, status string) (*[]models.Job, error) {
	headers := httpclient.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod, j.Settings.UsersID)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs?status=%s", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, svcID, status), headers)
	if err != nil {
		return nil, err
	}
	var jobs []models.Job
	err = httpclient.ConvertResp(resp, statusCode, &jobs)
	if err != nil {
		return nil, err
	}
	return &jobs, nil
}

func (j *SJobs) RetrieveByType(svcID, jobType string, page, pageSize int) (*[]models.Job, error) {
	headers := httpclient.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod, j.Settings.UsersID)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs?type=%s&pageNumber=%d&pageSize=%d", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, svcID, jobType, page, pageSize), headers)
	if err != nil {
		return nil, err
	}
	var jobs []models.Job
	err = httpclient.ConvertResp(resp, statusCode, &jobs)
	if err != nil {
		return nil, err
	}
	return &jobs, nil
}

func (j *SJobs) RetrieveByTarget(svcID, target string, page, pageSize int) (*[]models.Job, error) {
	var res []models.Job
	jobs, err := j.RetrieveByType(svcID, "worker", page, pageSize)
	if err != nil {
		return nil, err
	}
	for _, j := range *jobs {
		if j.Target == target {
			res = append(res, j)
		}
	}
	return &res, nil
}
