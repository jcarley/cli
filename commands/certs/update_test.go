package certs

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/commands/ssl"
	"github.com/daticahealth/cli/test"
)

var certUpdateTests = []struct {
	name        string
	pubKeyPath  string
	privKeyPath string
	selfSigned  bool
	resolve     bool
	expectErr   bool
}{
	{certName, pubKeyPath, privKeyPath, true, false, false},
	{certName, invalidPath, privKeyPath, true, false, true}, // invalid cert path
	{certName, pubKeyPath, invalidPath, true, false, true},  // invalid key path
	{certName, pubKeyPath, privKeyPath, false, false, true}, // cert not signed by CA
	{certName, pubKeyPath, privKeyPath, true, true, false},
	{"bad-cert-name", pubKeyPath, privKeyPath, true, false, true},
}

func TestCertsUpdate(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/certs/"+certName,
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "PUT")
			fmt.Fprint(w, fmt.Sprintf(`{"name":"%s"}`, certName))
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","label":"service_proxy"}]`, test.SvcID))
		},
	)

	for _, data := range certUpdateTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdUpdate(data.name, data.pubKeyPath, data.privKeyPath, data.selfSigned, data.resolve, New(settings), services.New(settings), ssl.New(settings))

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}
