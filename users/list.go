package users

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
)

func CmdList(myUsersID string, iu IUsers) error {
	envUsers, err := iu.List()
	if err != nil {
		return err
	}
	for _, userID := range envUsers.Users {
		if userID == myUsersID {
			logrus.Printf("%s (you)", userID)
		} else {
			defer logrus.Printf("%s", userID)
		}
	}
	return nil
}

func (u *SUsers) List() (*models.EnvironmentUsers, error) {
	headers := httpclient.GetHeaders(u.Settings.SessionToken, u.Settings.Version, u.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/users", u.Settings.PaasHost, u.Settings.PaasHostVersion, u.Settings.EnvironmentID), headers)
	if err != nil {
		return nil, err
	}
	var users models.EnvironmentUsers
	err = httpclient.ConvertResp(resp, statusCode, &users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}
