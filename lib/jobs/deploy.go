package jobs

import (
	"fmt"
	"strings"
)

func (j *SJobs) DeployRelease(releaseName, svcID string) error {
	return j.Deploy(true, releaseName, "", svcID)
}

func (j *SJobs) DeployTarget(target, svcID string) error {
	return j.Deploy(false, "", target, svcID)
}

func (j *SJobs) Redeploy(svcID string) error {
	return j.Deploy(true, "", "", svcID)
}

func (j *SJobs) Deploy(redeploy bool, releaseName, target, svcID string) error {
	var params = []string{}
	if releaseName != "" {
		params = append(params, fmt.Sprintf("release=%s", releaseName))
	}
	if redeploy {
		params = append(params, "redeploy=true")
	}
	if target != "" {
		params = append(params, fmt.Sprintf("target=%s", target))
	}
	headers := j.Settings.HTTPManager.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod, j.Settings.UsersID)
	resp, statusCode, err := j.Settings.HTTPManager.Post(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/deploy?%s", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, svcID, strings.Join(params, "&")), headers)
	if err != nil {
		return err
	}
	return j.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
