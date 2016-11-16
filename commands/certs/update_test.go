package certs

import (
	"testing"

	"github.com/catalyzeio/cli/test"
)

const (
	certsUpdateCommandName    = "certs"
	certsUpdateSubcommandName = "update"
	certsUpdateStandardOutput = `Updated 'example.com'
To make your updated cert go live, you must redeploy your service proxy with the "catalyze redeploy service_proxy" command
`
)

var certUpdateTests = []struct {
	env            string
	name           string
	pubKeyPath     string
	privKeyPath    string
	selfSigned     bool
	skipResolve    bool
	expectErr      bool
	expectedOutput string
}{
	{test.Alias, certName, pubKeyPath, privKeyPath, true, false, false, certsUpdateStandardOutput},
	{test.Alias, certName, invalidPath, privKeyPath, true, false, true, "\033[31m\033[1m[fatal] \033[0mA cert does not exist at path 'invalid-file.pem'\n"},
	{test.Alias, certName, pubKeyPath, invalidPath, true, false, true, "\033[31m\033[1m[fatal] \033[0mA private key does not exist at path 'invalid-file.pem'\n"},
	{test.Alias, certName, pubKeyPath, privKeyPath, false, false, false, "Incomplete certificate chain found, attempting to resolve this\n" + certsUpdateStandardOutput},
	{test.Alias, certName, pubKeyPath, privKeyPath, true, true, false, certsUpdateStandardOutput},
	{"bad-env", certName, pubKeyPath, privKeyPath, true, false, true, "\033[31m\033[1m[fatal] \033[0mNo environment named \"bad-env\" has been associated. Run \"catalyze associated\" to see what environments have been associated or run \"catalyze associate\" from a local git repo to create a new association\n"},
	{test.Alias, "bad-cert-name", pubKeyPath, privKeyPath, true, false, true, "\033[31m\033[1m[fatal] \033[0m(92002) Cert Not Found: Could not find a cert with the given name associated with the service.\n"},
}

func TestCertsUpdate(t *testing.T) {
	if err := test.SetUpGitRepo(); err != nil {
		t.Error(err)
		return
	}
	if err := test.SetUpAssociation(); err != nil {
		t.Error(err)
		return
	}
	if output, err := test.RunCommand(test.BinaryName, []string{"-E", test.Alias, certsCreateCommandName, certsCreateSubcommandName, certName, pubKeyPath, privKeyPath}); err != nil {
		t.Errorf("Error creating cert: %s - %s", err, output)
		return
	}

	for _, data := range certUpdateTests {
		t.Logf("Data: %+v", data)
		args := []string{"-E", data.env, certsUpdateCommandName, certsUpdateSubcommandName}
		if len(data.name) != 0 {
			args = append(args, data.name)
		}
		if len(data.pubKeyPath) != 0 {
			args = append(args, data.pubKeyPath)
		}
		if len(data.privKeyPath) != 0 {
			args = append(args, data.privKeyPath)
		}
		if data.selfSigned {
			args = append(args, "-s")
		}
		if data.skipResolve {
			args = append(args, "-r=false")
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
	if output, err := test.RunCommand(test.BinaryName, []string{"-E", test.Alias, certsRmCommandName, certsRmSubcommandName, certName}); err != nil {
		t.Errorf("Error removing cert: %s - %s", err, output)
		return
	}
}

func TestCertsUpdateNoAssociation(t *testing.T) {
	if err := test.ClearAssociations(); err != nil {
		t.Error(err)
		return
	}
	output, err := test.RunCommand(test.BinaryName, []string{certsUpdateCommandName, certsUpdateSubcommandName, certName, pubKeyPath, privKeyPath})
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
