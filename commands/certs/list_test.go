package certs

import (
	"testing"

	"github.com/catalyzeio/cli/test"
)

const (
	certsListCommandName    = "certs"
	certsListSubcommandName = "list"
	certsListStandardOutput = `NAME
sbox0513063.catalyzeapps.com
`
)

var certListTests = []struct {
	env            string
	expectErr      bool
	expectedOutput string
}{
	{test.Alias, false, certsListStandardOutput},
	{"bad-env", true, "\033[31m\033[1m[fatal] \033[0mNo environment named \"bad-env\" has been associated. Run \"catalyze associated\" to see what environments have been associated or run \"catalyze associate\" from a local git repo to create a new association\n"},
}

func TestCertsList(t *testing.T) {
	if err := test.SetUpGitRepo(); err != nil {
		t.Error(err)
		return
	}
	if err := test.SetUpAssociation(); err != nil {
		t.Error(err)
		return
	}

	for _, data := range certListTests {
		t.Logf("Data: %+v", data)
		args := []string{"-E", data.env}
		args = append(args, certsListCommandName, certsListSubcommandName)
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

func TestCertsListNoAssociation(t *testing.T) {
	if err := test.ClearAssociations(); err != nil {
		t.Error(err)
		return
	}
	output, err := test.RunCommand(test.BinaryName, []string{certsListCommandName, certsListSubcommandName})
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
