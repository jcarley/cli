package helpers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
)

// SignIn signs in the user and retrieves a session. The passed in Settings
// object is updated with the most up to date credentials
func SignIn(settings *models.Settings) {
	// if we're already signed in with a valid session, don't sign in again
	if verify(settings) {
		return
	}
	if settings.Username == "" || settings.Password == "" {
		promptForCredentials(settings)
	}
	login := models.Login{
		Username: settings.Username,
		Password: settings.Password,
	}
	b, err := json.Marshal(login)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	resp := httpclient.Post(b, fmt.Sprintf("%s%s/auth/signin", settings.BaasHost, config.BaasHostVersion), true, settings)
	var user models.User
	json.Unmarshal(resp, &user)
	settings.SessionToken = user.SessionToken
	settings.UsersID = user.UsersID
	config.SaveSettings(settings)
}

// verify tests whether or not the given session token is still valid
func verify(settings *models.Settings) bool {
	resp := httpclient.Get(fmt.Sprintf("%s%s/auth/verify", settings.BaasHost, config.BaasHostVersion), false, settings)
	m := make(map[string]string)
	json.Unmarshal(resp, &m)
	// the verify route returns userId and not usersId like everything else...
	if m["userId"] != "" {
		settings.UsersID = m["userId"]
	}
	return m["userId"] != ""
}
