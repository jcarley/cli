package associate

import (
	"testing"

	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/git"
	"github.com/catalyzeio/cli/models"
)

const validEnvName = "env"
const invalidEnvName = "badEnv"

const validAppName = "app01"
const invalidAppName = "badApp"

var associateTests = []struct {
	envLabel      string
	svcLabel      string
	alias         string
	remote        string
	defaultEnv    bool
	createGitRepo bool
	expectErr     bool
}{
	{validEnvName, validAppName, "e", "", false, true, false},
	{validEnvName, validAppName, "e", "", true, true, false},
	{validEnvName, validAppName, "e", "ctlyz", false, true, false},
	{validEnvName, validAppName, "", "", false, true, false},
	{validEnvName, invalidAppName, "e", "", false, true, true},
	{invalidEnvName, validAppName, "e", "", false, true, true},
	{validEnvName, validAppName, "e", "", false, false, true},
}

func TestAssociate(t *testing.T) {
	for _, data := range associateTests {
		t.Logf("%+v\n", data)
		createGitRepo(data.envLabel, data.createGitRepo)

		settings := getSettings()
		sa := SAssociate{
			Settings: settings,
			Git: &git.MGit{
				ReturnError: false,
			},
			Environments: &environments.MEnvironments{
				ReturnError: false,
			},
			EnvLabel:   data.envLabel,
			SvcLabel:   data.svcLabel,
			Alias:      data.alias,
			Remote:     data.remote,
			DefaultEnv: data.defaultEnv,
		}
		err := sa.Associate()
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s\n", err.Error())
		}
		name := data.alias
		if name == "" {
			name = data.envLabel
		}

		found := false
		for _, env := range settings.Environments {
			if env.Name == name {
				found = true
				break
			}
		}
		if !found {
			t.Error("Environment not added to the settings list of environments or the alias was not used")
		}
		if data.defaultEnv && settings.Default != name {
			t.Error("Default environment specified but was not stored in the settings")
		}

		expectedRemote := data.remote
		if expectedRemote == "" {
			expectedRemote = "catalyze"
		}
		remotes, err := git.New().List()
		if err != nil {
			t.Errorf("Error listing git remotes: %s\n", err.Error())
		}
		found = false
		for _, r := range remotes {
			if r == data.remote {
				found = true
			}
		}
		if !found {
			t.Errorf("Proper git remote not listed. Found '%+v' instead\n", remotes)
		}

		destroyGitRepo(data.envLabel)
	}
}

func getSettings() *models.Settings {
	return &models.Settings{
		BaasHost:        config.BaasHost,
		PaasHost:        config.PaasHost,
		Username:        "test",
		Password:        "test",
		EnvironmentID:   "1234",
		ServiceID:       "5678",
		Pod:             "pod01",
		EnvironmentName: validEnvName,
		SessionToken:    "1234567890",
		UsersID:         "0987654321",
		Environments:    make(map[string]models.AssociatedEnv, 0),
		Default:         "",
		Pods:            &[]models.Pod{},
	}
}

func createGitRepo(name string, create bool) {

}

func destroyGitRepo(name string) {

}
