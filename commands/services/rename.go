package services

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/httpclient"
)

func CmdRename(svcID, label string, is IServices) error {
	data := map[string]string{}
	data["label"] = label
	err := is.Update(svcID, data)
	if err != nil {
		return err
	}
	logrus.Printf("Successfully renamed your service to %s", label)
	return nil
}

func (s *SServices) Update(svcID string, updates map[string]string) error {
	b, err := json.Marshal(updates)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(s.Settings.SessionToken, s.Settings.Version, s.Settings.Pod, s.Settings.UsersID)
	resp, statusCode, err := httpclient.Put(b, fmt.Sprintf("%s%s/environments/%s/services/%s", s.Settings.PaasHost, s.Settings.PaasHostVersion, s.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
