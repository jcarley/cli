package users

import (
	"encoding/json"
	"fmt"

	"github.com/catalyzeio/cli/httpclient"
)

func CmdAdd(usersID string, iu IUsers) error {
	fmt.Println("WARNING: This command is deprecated. Please use \"catalyze invites send\" instead.")
	err := iu.Add(usersID)
	if err != nil {
		return err
	}
	fmt.Println("Added.")
	return nil
}

func (u *SUsers) Add(usersID string) error {
	body := map[string]string{}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(u.Settings.SessionToken, u.Settings.Version, u.Settings.Pod)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/users/%s", u.Settings.PaasHost, u.Settings.PaasHostVersion, u.Settings.EnvironmentID, usersID), headers)
	if err != nil {
		return nil
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
