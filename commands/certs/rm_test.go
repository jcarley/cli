package certs

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/test"
)

var certRmTests = []struct {
	name       string
	downStream string
	expectErr  bool
}{
	{certName, test.DownStream, false},
	{"bad-cert-name", test.DownStream, true},
}

func TestCertsRm(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/certs/"+certName,
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "DELETE")
			fmt.Fprint(w, "")
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","label":"%s"}]`, test.SvcID, test.DownStream))
		},
	)

	for _, data := range certRmTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdRm(data.name, New(settings), services.New(settings), data.downStream)

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}
