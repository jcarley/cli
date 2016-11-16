package clear

import (
	"testing"

	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/test"
)

const (
	commandName    = "clear"
	standardOutput = "All specified settings have been cleared\n"
)

var settings = models.Settings{
	OrgID:          "org1234",
	PrivateKeyPath: "",
	SessionToken:   "1234",
	UsersID:        "user1234",
	Environments: map[string]models.AssociatedEnv{
		"ctest": models.AssociatedEnv{
			EnvironmentID: "env1234",
			ServiceID:     "svc1234",
			Directory:     "/dir",
			Name:          "test",
			Pod:           "air-force-pod",
			OrgID:         "org1234",
		},
	},
	Default: "ctest",
	Pods: &[]models.Pod{
		{
			Name:                 "air-force-pod",
			PHISafe:              true,
			ImportRequiresLength: true,
		},
	},
	PodCheck: 1478899890,
}

var clearTests = []struct {
	env            string
	privKey        bool
	session        bool
	envs           bool
	defaultEnv     bool
	pods           bool
	all            bool
	expectErr      bool
	expectedOutput string
}{
	{test.Alias, false, false, false, false, false, false, false, "No settings were specified. To see available options, run \"catalyze clear --help\"\n"},
	{test.Alias, true, false, false, false, false, false, false, standardOutput},
	{test.Alias, false, true, false, false, false, false, false, standardOutput},
	{test.Alias, false, false, true, false, false, false, false, standardOutput},
	{test.Alias, false, false, false, true, false, false, false, "\033[33m\033[1m[warning] \033[0mThe \"--default\" flag has been deprecated! It will be removed in a future version.\n" + standardOutput},
	{test.Alias, false, false, false, false, true, false, false, standardOutput},
	{test.Alias, false, false, false, false, false, true, false, "\033[33m\033[1m[warning] \033[0mThe \"--default\" flag has been deprecated! It will be removed in a future version.\n" + standardOutput},
	{"bad-env", true, false, false, false, false, false, true, "\033[31m\033[1m[fatal] \033[0mNo environment named \"bad-env\" has been associated. Run \"catalyze associated\" to see what environments have been associated or run \"catalyze associate\" from a local git repo to create a new association\n"},
}

func TestClear(t *testing.T) {
	for _, data := range clearTests {
		t.Logf("Data: %+v", data)
		if err := restoreSettings(); err != nil {
			t.Error(err)
			return
		}
		args := []string{"-E", data.env, commandName}
		if data.privKey {
			args = append(args, "--private-key")
		}
		if data.session {
			args = append(args, "--session")
		}
		if data.envs {
			args = append(args, "--environments")
		}
		if data.defaultEnv {
			args = append(args, "--default")
		}
		if data.pods {
			args = append(args, "--pods")
		}
		if data.all {
			args = append(args, "--all")
		}
		output, err := test.RunCommand(test.BinaryName, args)
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", output)
			continue
		}
		if output != data.expectedOutput {
			t.Errorf("Expected: %s. Found: %s", data.expectedOutput, output)
			continue
		}
	}
}

func TestClearNoAssociation(t *testing.T) {
	if err := test.ClearAssociations(); err != nil {
		t.Error(err)
		return
	}
	output, err := test.RunCommand(test.BinaryName, []string{commandName, "--all"})
	if err != nil {
		t.Errorf("Unexpected error : %s - %s", err, output)
		return
	}
	expectedOutput := "\033[33m\033[1m[warning] \033[0mThe \"--default\" flag has been deprecated! It will be removed in a future version.\n" + standardOutput
	if output != expectedOutput {
		t.Errorf("Expected: %s. Found: %s", expectedOutput, output)
		return
	}
}

func restoreSettings() error {
	config.SaveSettings(&settings)
	return nil
}
