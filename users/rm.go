package users

import (
	"fmt"

	"github.com/catalyzeio/cli/httpclient"
)

func CmdRm(usersID string, iu IUsers) error {
	err := iu.Rm(usersID)
	if err != nil {
		return err
	}
	fmt.Println("Removed.")
	return nil
}

func (u *SUsers) Rm(usersID string) error {
	headers := httpclient.GetHeaders(u.Settings.SessionToken, u.Settings.Version, u.Settings.Pod)
	resp, statusCode, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/environments/%s/users/%s", u.Settings.PaasHost, u.Settings.PaasHostVersion, u.Settings.EnvironmentID, usersID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
