package auth

import (
	"encoding/json"
	"fmt"

	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
)

// Signin signs in a user and returns the representative user model. If an
// error occurs, nil is returned for the user and the error field is populated.
func (a *SAuth) Signin() (*models.User, error) {
	// if we're already signed in with a valid session, don't sign in again
	if a.Verify() == nil {
		return &models.User{
			Username:     "",
			SessionToken: a.Settings.SessionToken,
			UsersID:      a.Settings.UsersID,
		}, nil
	}
	//var username, password string
	login := models.Login{
		Identifier: a.Settings.Username,
		Password:   a.Settings.Password,
	}
	if a.Settings.Username == "" || a.Settings.Password == "" {
		username, password, err := a.Prompts.UsernamePassword()
		if err != nil {
			return nil, err
		}
		login = models.Login{
			Identifier: username,
			Password:   password,
		}
	}

	b, err := json.Marshal(login)
	if err != nil {
		return nil, err
	}
	headers := httpclient.GetHeaders(a.Settings.SessionToken, a.Settings.Version, a.Settings.Pod)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/auth/signin", a.Settings.AuthHost, a.Settings.AuthHostVersion), headers)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = httpclient.ConvertResp(resp, statusCode, &user)
	if err != nil {
		return nil, err
	}
	a.Settings.SessionToken = user.SessionToken
	a.Settings.UsersID = user.UsersID
	return &user, nil
}

// Signout signs out a user by their session token.
func (a *SAuth) Signout() error {
	headers := httpclient.GetHeaders(a.Settings.SessionToken, a.Settings.Version, a.Settings.Pod)
	resp, statusCode, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/auth/signout", a.Settings.AuthHost, a.Settings.AuthHostVersion), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}

// Verify verifies if a given session token is still valid or not. If it is
// valid, the returned error will be nil.
func (a *SAuth) Verify() error {
	headers := httpclient.GetHeaders(a.Settings.SessionToken, a.Settings.Version, a.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/auth/verify", a.Settings.AuthHost, a.Settings.AuthHostVersion), headers)
	if err != nil {
		return err
	}
	m := make(map[string]string)
	err = httpclient.ConvertResp(resp, statusCode, &m)
	if err != nil {
		return err
	}
	// the verify route returns userId and not usersId like everything else...
	if m["id"] != "" {
		a.Settings.UsersID = m["id"]
		return nil
	}
	return fmt.Errorf("Invalid session token: %s", string(resp))
}
