package environments

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
)

// CmdEnvironments lists all environments which the user has access to
func CmdEnvironments(environments IEnvironments) error {
	envs, err := environments.List()
	if err != nil {
		return err
	}
	for _, env := range *envs {
		logrus.Printf("%+v", env)
		//logrus.Printf("%s: %s", env.Data.Name, env.ID)
	}
	if len(*envs) == 0 {
		logrus.Println("no environments found")
	}
	return nil
}

func (e *SEnvironments) List() (*[]models.Environment, error) {
	var allEnvs []models.Environment
	logrus.Debugf("pods %+v", e.Settings.Pods)
	for _, pod := range *e.Settings.Pods {
		headers := httpclient.GetHeaders(e.Settings.SessionToken, e.Settings.Version, pod.Name)
		resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments", e.Settings.PaasHost, e.Settings.PaasHostVersion), headers)
		if err != nil {
			return nil, err
		}
		var envs []models.Environment
		err = httpclient.ConvertResp(resp, statusCode, &envs)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(envs); i++ {
			envs[i].Pod = pod.Name
			logrus.Debugf("e %+v", envs[i])
		}
		allEnvs = append(allEnvs, envs...)
	}
	logrus.Debugf("all envs %+v", allEnvs)
	return &allEnvs, nil
}

func (e *SEnvironments) Retrieve(envID string) (*models.Environment, error) {
	headers := httpclient.GetHeaders(e.Settings.SessionToken, e.Settings.Version, e.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s", e.Settings.PaasHost, e.Settings.PaasHostVersion, envID), headers)
	if err != nil {
		return nil, err
	}
	var env models.Environment
	err = httpclient.ConvertResp(resp, statusCode, &env)
	if err != nil {
		return nil, err
	}
	return &env, nil
}
