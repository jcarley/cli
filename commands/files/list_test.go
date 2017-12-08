package files

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/test"
)

var listTests = []struct {
	svcName        string
	showTimestamps bool
	expectErr      bool
}{
	{test.SvcLabel, false, false},
	{test.SvcLabel, true, false},
	{"invalid-svc", false, true},
}

func TestList(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","label":"%s"}]`, test.SvcID, test.SvcLabel))
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/files",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":1,"name":"%s","created_at":"%s","updated_at":"%s"}]`, fileName, "2016-11-16T16:31:12", "2017-11-16T16:31:12"))
		},
	)

	for _, data := range listTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdList(data.svcName, data.showTimestamps, New(settings), services.New(settings))

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}
