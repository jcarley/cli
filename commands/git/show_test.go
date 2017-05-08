package git

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/test"
)

var showTests = []struct {
	svcLabel  string
	expectErr bool
}{
	{test.SvcLabel, false},
	{"invalid-svc", true},
}

func TestShow(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprintf(w, `[{"id":"%s","label":"%s","source":"git@github.com/github/github.git"}]`, test.SvcID, test.SvcLabel)
		},
	)

	for _, data := range showTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdShow(data.svcLabel, services.New(settings))

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}
