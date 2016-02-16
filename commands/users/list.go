package users

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func CmdList(myUsersID string, iu IUsers) error {
	orgUsers, err := iu.List()
	if err != nil {
		return err
	}
	for _, user := range *orgUsers {
		if user.ID == myUsersID {
			logrus.Printf("%s (you)", user.ID)
		} else {
			defer logrus.Printf("%s", user.ID)
		}
	}
	return nil
}

func (u *SUsers) List() (*[]models.OrgUser, error) {
	headers := httpclient.GetHeaders(u.Settings.SessionToken, u.Settings.Version, u.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/orgs/%s/users", u.Settings.AuthHost, u.Settings.AuthHostVersion, u.Settings.OrgID), headers)
	if err != nil {
		return nil, err
	}
	var users []models.OrgUser
	err = httpclient.ConvertResp(resp, statusCode, &users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}
