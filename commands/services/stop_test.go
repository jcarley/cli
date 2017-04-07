package services

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/daticahealth/cli/lib/jobs"
	"github.com/daticahealth/cli/test"
)

var servicesStopTests = []struct {
	svcName   string
	expectErr bool
}{
	{test.SvcLabel, false},
	{test.SvcLabelAlt, true},
	{"database-1", true},
	{"invalid-svc-name", true},
}

func TestServicesStop(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","label":"%s","name":"code","redeployable":false},{"id":"%s","label":"%s","name":"code","redeployable":true},{"id":"%s","label":"database-1","name":"postgresql","redeployable":true}]`, test.SvcIDAlt, test.SvcLabelAlt, test.SvcID, test.SvcLabel, "3"))
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/jobs",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, `[{"id":"1"}]`)
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/jobs/1",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "DELETE")
			w.WriteHeader(204)
		},
	)

	for _, data := range servicesStopTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdStop(data.svcName, New(settings), jobs.New(settings), &test.FakePrompts{})

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}
