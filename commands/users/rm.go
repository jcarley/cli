package users

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

func CmdRm(email string, iu IUsers) error {
	orgUsers, err := iu.List()
	if err != nil {
		return err
	}
	usersID := ""
	for _, u := range *orgUsers {
		if u.Email == email {
			usersID = u.ID
			break
		}
	}
	if usersID == "" {
		return fmt.Errorf("A user with email %s was not found", email)
	}

	err = iu.Rm(usersID)
	if err != nil {
		return err
	}
	logrus.Printf("Removed %s from your environment's organization.", email)
	return nil
}

func (u *SUsers) Rm(usersID string) error {
	headers := u.Settings.HTTPManager.GetHeaders(u.Settings.SessionToken, u.Settings.Version, u.Settings.Pod, u.Settings.UsersID)
	resp, statusCode, err := u.Settings.HTTPManager.Delete(nil, fmt.Sprintf("%s%s/orgs/%s/users/%s", u.Settings.AuthHost, u.Settings.AuthHostVersion, u.Settings.OrgID, usersID), headers)
	if err != nil {
		return err
	}
	return u.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
