package jobs

import "fmt"

func (j *SJobs) Delete(jobID, svcID string) error {
	headers := j.Settings.HTTPManager.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod, j.Settings.UsersID)
	resp, statusCode, err := j.Settings.HTTPManager.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs/%s", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, svcID, jobID), headers)
	if err != nil {
		return err
	}
	return j.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
