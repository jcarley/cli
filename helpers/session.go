package helpers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/catalyzeio/catalyze/config"
	"github.com/catalyzeio/catalyze/httpclient"
	"github.com/catalyzeio/catalyze/models"
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

func promptForCredentials(settings *models.Settings) {
	var username string
	fmt.Print("Username: ")
	in := bufio.NewReader(os.Stdin)
	username, err := in.ReadString('\n')
	if err != nil {
		panic(errors.New("Invalid username"))
	}
	username = strings.TrimRight(username, "\n")
	if runtime.GOOS == "windows" {
		username = strings.TrimRight(username, "\r")
	}
	settings.Username = username
	fmt.Print("Password: ")
	bytes, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println("")
	settings.Password = string(bytes)
}
