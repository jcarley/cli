package environments

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/daticahealth/cli/test"
)

var renameTests = []struct {
	envID     string
	name      string
	expectErr bool
}{
	{test.EnvID, test.EnvName, false},
	{test.EnvIDAlt, test.EnvNameAlt, true},
}

func TestRename(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID,
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "PUT")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","name":"%s","namespace":"%s","organizationId":"%s"}]`, test.EnvID, test.EnvName, test.Namespace, test.OrgID))
		},
	)

	for _, data := range renameTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdRename(data.envID, data.name, New(settings))

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}
