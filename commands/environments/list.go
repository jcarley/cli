package environments

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

// CmdList lists all environments which the user has access to
func CmdList(environments IEnvironments) error {
	envs, errs := environments.List()
	if envs == nil || len(*envs) == 0 {
		logrus.Println("no environments found")
	} else {
		for _, env := range *envs {
			logrus.Printf("%s: %s", env.Name, env.ID)
		}
	}
	if errs != nil && len(errs) > 0 {
		for pod, err := range errs {
			defer logrus.Errorf("Failed to list environments for pod \"%s\": %s", pod, err)
		}
		return errors.New("Failures listing environments. Please try again or contact support@catalyze.io if this issue persists.")
	}
	return nil
}

func (e *SEnvironments) List() (*[]models.Environment, map[string]error) {
	allEnvs := []models.Environment{}
	errs := map[string]error{}
	for _, pod := range *e.Settings.Pods {
		headers := httpclient.GetHeaders(e.Settings.SessionToken, e.Settings.Version, pod.Name, e.Settings.UsersID)
		resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments", e.Settings.PaasHost, e.Settings.PaasHostVersion), headers)
		if err != nil {
			errs[pod.Name] = err
			continue
		}
		var envs []models.Environment
		err = httpclient.ConvertResp(resp, statusCode, &envs)
		if err != nil {
			errs[pod.Name] = err
			continue
		}
		for i := 0; i < len(envs); i++ {
			envs[i].Pod = pod.Name
		}
		allEnvs = append(allEnvs, envs...)
	}
	return &allEnvs, errs
}

func (e *SEnvironments) Retrieve(envID string) (*models.Environment, error) {
	headers := httpclient.GetHeaders(e.Settings.SessionToken, e.Settings.Version, e.Settings.Pod, e.Settings.UsersID)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s", e.Settings.PaasHost, e.Settings.PaasHostVersion, envID), headers)
	if err != nil {
		return nil, err
	}
	var env models.Environment
	err = httpclient.ConvertResp(resp, statusCode, &env)
	if err != nil {
		return nil, err
	}
	env.Pod = e.Settings.Pod
	return &env, nil
}
