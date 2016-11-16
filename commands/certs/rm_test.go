package certs

import (
	"testing"

	"github.com/catalyzeio/cli/test"
)

const (
	certsRmCommandName    = "certs"
	certsRmSubcommandName = "rm"
	certsRmStandardOutput = "Removed 'example.com'\n"
)

var certRmTests = []struct {
	env            string
	name           string
	expectErr      bool
	expectedOutput string
}{
	{test.Alias, certName, false, certsRmStandardOutput},
	{"bad-env", certName, true, "\033[31m\033[1m[fatal] \033[0mNo environment named \"bad-env\" has been associated. Run \"catalyze associated\" to see what environments have been associated or run \"catalyze associate\" from a local git repo to create a new association\n"},
	{test.Alias, "bad-cert-name", true, "\033[31m\033[1m[fatal] \033[0m(92002) Cert Not Found: Could not find a cert with the given name associated with the service.\n"},
}

func TestCertsRm(t *testing.T) {
	if err := test.SetUpGitRepo(); err != nil {
		t.Error(err)
		return
	}
	if err := test.SetUpAssociation(); err != nil {
		t.Error(err)
		return
	}

	for _, data := range certRmTests {
		t.Logf("Data: %+v", data)
		if !data.expectErr {
			if output, err := test.RunCommand(test.BinaryName, []string{"-E", test.Alias, certsCreateCommandName, certsCreateSubcommandName, certName, pubKeyPath, privKeyPath}); err != nil {
				t.Errorf("Error creating cert: %s - %s", err, output)
				continue
			}
		}
		args := []string{"-E", data.env, certsRmCommandName, certsRmSubcommandName}
		if len(data.name) != 0 {
			args = append(args, data.name)
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

func TestCertsRmNoAssociation(t *testing.T) {
	if err := test.ClearAssociations(); err != nil {
		t.Error(err)
		return
	}
	output, err := test.RunCommand(test.BinaryName, []string{certsRmCommandName, certsRmSubcommandName, certName})
	if err == nil {
		t.Errorf("Expected error but no error returned: %s", output)
		return
	}
	expectedOutput := "\033[31m\033[1m[fatal] \033[0mNo Catalyze environment has been associated. Run \"catalyze associate\" from a local git repo first\n"
	if output != expectedOutput {
		t.Errorf("Expected: %s. Found: %s", expectedOutput, output)
		return
	}
}
