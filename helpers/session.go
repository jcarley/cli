package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/catalyzeio/catalyze/config"
	"github.com/catalyzeio/catalyze/httpclient"
	"github.com/catalyzeio/catalyze/models"
	"github.com/docker/docker/pkg/term"
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
	resp := httpclient.Post(b, fmt.Sprintf("%s/v2/auth/signin", settings.BaasHost), true, settings)
	var user models.User
	json.Unmarshal(resp, &user)
	settings.SessionToken = user.SessionToken
	settings.UsersID = user.UsersID
	config.SaveSettings(settings)
}

// verify tests whether or not the given session token is still valid
func verify(settings *models.Settings) bool {
	resp := httpclient.Get(fmt.Sprintf("%s/v2/auth/verify", settings.BaasHost), false, settings)
	m := make(map[string]string)
	json.Unmarshal(resp, &m)
	// the verify route returns userId and not usersId like everything else...
	if m["userId"] != "" {
		settings.UsersID = m["userId"]
	}
	return m["userId"] != ""
}

func promptForCredentials(settings *models.Settings) {
	var username string
	fmt.Print("Username: ")
	fmt.Scanln(&username)
	settings.Username = username
	fmt.Print("Password: ")
	var fd uintptr
	if runtime.GOOS == "windows" {
		stdIn, _, _ := term.StdStreams()
		fd, _ = term.GetFdInfo(stdIn)
	}
	bytes, _ := terminal.ReadPassword(int(fd))
	fmt.Println("")
	settings.Password = string(bytes)
}
