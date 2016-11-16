package associate

import (
	"strings"
	"testing"

	"github.com/catalyzeio/cli/test"
)

const (
	commandName    = "associate"
	standardOutput = `Existing git remotes named "catalyze" will be overwritten
"catalyze" remote added.
Your git repository "catalyze" has been associated with code service "code-1" and environment "ctest"
After associating to an environment, you need to add a cert with the "catalyze certs create" command, if you have not done so already
`
)

var associateTests = []struct {
	envLabel       string
	svcLabel       string
	alias          string
	remote         string
	defaultEnv     bool
	expectErr      bool
	expectedOutput string
}{
	{test.EnvLabel, test.SvcLabel, test.Alias, "", false, false, standardOutput},
	{"bad-env", test.SvcLabel, test.Alias, "", false, true, "Existing git remotes named \"catalyze\" will be overwritten\n\033[31m\033[1m[fatal] \033[0mNo environment with name \"bad-env\" found\n"},
	{test.EnvLabel, "bad-svc", test.Alias, "", false, true, "Existing git remotes named \"catalyze\" will be overwritten\n\033[31m\033[1m[fatal] \033[0mNo code service found with label \"bad-svc\". Code services found: code-1\n"},
	{test.EnvLabel, test.SvcLabel, "", "", false, false, "Existing git remotes named \"catalyze\" will be overwritten\n\"catalyze\" remote added.\nYour git repository \"catalyze\" has been associated with code service \"code-1\" and environment \"cli-integration-tests\"\nAfter associating to an environment, you need to add a cert with the \"catalyze certs create\" command, if you have not done so already\n"},
	{test.EnvLabel, test.SvcLabel, test.Alias, "cz", false, false, "Existing git remotes named \"cz\" will be overwritten\n\"cz\" remote added.\nYour git repository \"cz\" has been associated with code service \"code-1\" and environment \"ctest\"\nAfter associating to an environment, you need to add a cert with the \"catalyze certs create\" command, if you have not done so already\n"},
	{test.EnvLabel, test.SvcLabel, test.Alias, "", true, false, "\033[33m\033[1m[warning] \033[0mThe \"--default\" flag has been deprecated! It will be removed in a future version.\n" + standardOutput},
	{"", test.SvcLabel, test.Alias, "", false, true, "Existing git remotes named \"catalyze\" will be overwritten\n\033[31m\033[1m[fatal] \033[0mNo environment with name \"\" found\n"},
	{test.EnvLabel, "", test.Alias, "", false, true, "Existing git remotes named \"catalyze\" will be overwritten\n\033[31m\033[1m[fatal] \033[0mNo code service found with label \"\". Code services found: code-1\n"},
}

func TestAssociate(t *testing.T) {
	if err := test.SetUpGitRepo(); err != nil {
		t.Error(err)
		t.Fail()
	}
	for _, data := range associateTests {
		t.Logf("Data: %+v", data)
		args := []string{commandName, data.envLabel, data.svcLabel}
		if len(data.alias) != 0 {
			args = append(args, "-a", data.alias)
		}
		expectedRemote := "catalyze"
		if len(data.remote) != 0 {
			expectedRemote = data.remote
			args = append(args, "-r", data.remote)
		}
		if data.defaultEnv {
			args = append(args, "-d")
		}

		output, err := test.RunCommand(test.BinaryName, args)
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", output)
			continue
		}
		if output != data.expectedOutput {
			t.Errorf("Expected: %v. Found: %v", data.expectedOutput, output)
			continue
		}
		output, err = test.RunCommand("git", []string{"remote", "-v"})
		if err != nil {
			t.Errorf("Unexpected error running 'git remote -v': %s", output)
			continue
		}
		if !strings.Contains(output, expectedRemote) {
			t.Errorf("Git remote not added. Expected: %s. Found %s", expectedRemote, output)
			continue
		}
	}
}
