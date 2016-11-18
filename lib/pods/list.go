package pods

import (
	"fmt"

	"github.com/catalyzeio/cli/models"
)

func (p *SPods) List() (*[]models.Pod, error) {
	headers := p.Settings.HTTPManager.GetHeaders(p.Settings.SessionToken, p.Settings.Version, p.Settings.Pod, p.Settings.UsersID)
	resp, statusCode, err := p.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/pods", p.Settings.PaasHost, p.Settings.PaasHostVersion), headers)
	if err != nil {
		return nil, err
	}
	var podWrapper models.PodWrapper
	err = p.Settings.HTTPManager.ConvertResp(resp, statusCode, &podWrapper)
	if err != nil {
		return nil, err
	}
	return podWrapper.Pods, nil
}
