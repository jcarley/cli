package volumes

import (
	"fmt"

	"github.com/daticahealth/cli/models"
)

func (v *SVolumes) List(svcID string) (*[]models.Volume, error) {
	headers := v.Settings.HTTPManager.GetHeaders(v.Settings.SessionToken, v.Settings.Version, v.Settings.Pod, v.Settings.UsersID)
	resp, statusCode, err := v.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/volumes", v.Settings.PaasHost, v.Settings.PaasHostVersion, v.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return nil, err
	}
	var volumes []models.Volume
	err = v.Settings.HTTPManager.ConvertResp(resp, statusCode, &volumes)
	if err != nil {
		return nil, err
	}
	return &volumes, nil
}
