package services

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/daticahealth/cli/lib/jobs"
	"github.com/daticahealth/cli/lib/volumes"
	"github.com/daticahealth/cli/test"
)

var servicesStopTests = []struct {
	svcName   string
	pod       string
	expectErr bool
}{
	{test.SvcLabel, test.Pod, false},
	{test.SvcLabelAlt, test.Pod, true},
	{"database-1", test.Pod, true},
	{"invalid-svc-name", test.Pod, true},
	{"database-2", "csb01", true},
	{"database-2", test.Pod, true},
}

func TestServicesStop(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","label":"%s","name":"code","redeployable":false},{"id":"%s","label":"%s","name":"code","redeployable":true},{"id":"%s","label":"database-1","name":"postgresql","redeployable":true},{"id":"%s","label":"database-2","name":"postgresql","redeployable":true}]`, test.SvcIDAlt, test.SvcLabelAlt, test.SvcID, test.SvcLabel, "3", "4"))
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
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/volumes",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, `[]`)
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services/4/volumes",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, `[{"id":"v","type":"simple","size":"20"}]`)
		},
	)

	for _, data := range servicesStopTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdStop(data.svcName, data.pod, New(settings), jobs.New(settings), volumes.New(settings), &test.FakePrompts{})

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}
