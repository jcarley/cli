package auth

import (
	"encoding/json"
	"fmt"

	"github.com/catalyzeio/cli/config"
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
	var login models.Login
	if a.Settings.Username == "" || a.Settings.Password == "" {
		username, password, err := a.Prompts.UsernamePassword()
		if err != nil {
			return nil, err
		}
		login = models.Login{
			Username: username,
			Password: password,
		}
	}

	b, err := json.Marshal(login)
	if err != nil {
		// TODO this is not a nice error, fix it
		return nil, err
	}
	resp := httpclient.Post(b, fmt.Sprintf("%s%s/auth/signin", a.Settings.BaasHost, config.BaasHostVersion), true, a.Settings)
	var user models.User
	json.Unmarshal(resp, &user)
	// TODO settings is the absolute wrong place for constant changing data like session token and users ID
	a.Settings.SessionToken = user.SessionToken
	a.Settings.UsersID = user.UsersID
	config.SaveSettings(a.Settings)
	return &user, nil
}

// Signout signs out a user by their session token.
func (a *SAuth) Signout() error {
	return nil
}

// Verify verifies if a given session token is still valid or not. If it is
// valid, the returned error will be nil.
func (a *SAuth) Verify() error {
	resp := httpclient.Get(fmt.Sprintf("%s%s/auth/verify", a.Settings.BaasHost, config.BaasHostVersion), false, a.Settings)
	m := make(map[string]string)
	json.Unmarshal(resp, &m)
	// the verify route returns userId and not usersId like everything else...
	if m["userId"] != "" {
		a.Settings.UsersID = m["userId"]
		return nil
	}
	// TODO parse the resp in a proper error format
	return fmt.Errorf("Invalid session token: %s", string(resp))
}
