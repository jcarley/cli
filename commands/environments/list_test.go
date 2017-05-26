package environments

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/daticahealth/cli/test"
)

func TestList(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			if r.Header.Get("X-Pod-ID") == test.Pod {
				fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","name":"%s","namespace":"%s","organizationId":"%s"}]`, test.EnvID, test.EnvName, test.Namespace, test.OrgID))
			} else {
				fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","name":"%s","namespace":"%s","organizationId":"%s"}]`, test.EnvIDAlt, test.EnvNameAlt, test.NamespaceAlt, test.OrgIDAlt))
			}
		},
	)

	err := CmdList(New(settings))

	// assert
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestListWithPodError(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			w.WriteHeader(500)
		},
	)

	err := CmdList(New(settings))

	// assert
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}
