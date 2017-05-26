package disassociate

import (
	"testing"

	"github.com/daticahealth/cli/test"
)

var disassociateTests = []struct {
	name      string
	expectErr bool
}{
	{test.Alias, false},
	{"bad-alias", true},
}

func TestDisassociate(t *testing.T) {
	settings := test.GetSettings("")

	for _, data := range disassociateTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdDisassociate(data.name, New(settings))

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
		if _, present := settings.Environments[data.name]; present {
			t.Errorf("Environment not removed from settings")
			continue
		}
	}
}
