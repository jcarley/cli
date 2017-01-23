package environments

import (
	"encoding/json"
	"fmt"
)

func CmdRename(envID, name string, ie IEnvironments) error {
	data := map[string]string{}
	data["name"] = name
	return ie.Update(envID, data)
}

func (e *SEnvironments) Update(envID string, updates map[string]string) error {
	b, err := json.Marshal(updates)
	if err != nil {
		return err
	}
	headers := e.Settings.HTTPManager.GetHeaders(e.Settings.SessionToken, e.Settings.Version, e.Settings.Pod, e.Settings.UsersID)
	resp, statusCode, err := e.Settings.HTTPManager.Put(b, fmt.Sprintf("%s%s/environments/%s", e.Settings.PaasHost, e.Settings.PaasHostVersion, envID), headers)
	if err != nil {
		return err
	}
	return e.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
