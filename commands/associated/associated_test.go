package associated

import (
	"regexp"
	"testing"

	"github.com/catalyzeio/cli/test"
)

const (
	commandName    = "associated"
	standardOutput = `ctest:
\s+Environment ID:[\s]+[^\s]+
\s+Environment Name: cli-integration-tests
\s+Service ID:[\s]+[^\s]+
\s+Associated at:[\s]+[^\s]+
\s+Pod:[\s]+[^\s]+
\s+Organization ID:[\s]+[^\s]+`
)

var associatedTests = []struct {
	associateFirst bool
	expectErr      bool
	expectedOutput string
}{
	{true, false, standardOutput},
	{false, false, "No environments have been associated"},
}

func TestAssociated(t *testing.T) {
	if err := test.SetUpGitRepo(); err != nil {
		t.Error(err)
		t.Fail()
	}
	for _, data := range associatedTests {
		t.Logf("Data: %+v", data)
		if err := test.ClearAssociations(); err != nil {
			t.Error(err)
			continue
		}
		if data.associateFirst {
			if err := test.SetUpAssociation(); err != nil {
				t.Error(err)
				continue
			}
		}
		r := regexp.MustCompile(data.expectedOutput)
		output, err := test.RunCommand(test.BinaryName, []string{commandName})
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", output)
			continue
		}
		if !r.MatchString(output) {
			t.Errorf("Expected: %s. Found: %s", data.expectedOutput, output)
			continue
		}
	}
}
