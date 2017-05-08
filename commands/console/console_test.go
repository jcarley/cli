package console

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/lib/jobs"
	"github.com/daticahealth/cli/test"
)

var consoleTests = []struct {
	svcName   string
	command   string
	expectErr bool
}{
	{test.SvcLabel, "echo 1", true}, // stdin is not a terminal
	{"invalid-svc", "echo 1", true},
}

func TestConsole(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","label":"%s"}]`, test.SvcID, test.SvcLabel))
		},
	)

	for _, data := range consoleTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdConsole(data.svcName, data.command, New(settings, jobs.New(settings)), services.New(settings))

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}
