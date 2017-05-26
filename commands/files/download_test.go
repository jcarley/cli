package files

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/test"
)

const fileContents = "text"

var downloadTests = []struct {
	svcName   string
	fileName  string
	output    string
	force     bool
	expectErr bool
}{
	{test.SvcLabel, fileName, "", false, false},
	{test.SvcLabel, fileName, "", true, false},
	{test.SvcLabel, fileName, "output.txt", true, false},
	{test.SvcLabel, "invalid-file", "", false, true},
	{"invalid-svc", fileName, "", false, true},
}

func TestDownload(t *testing.T) {
	os.Remove("output.txt")
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
			fmt.Fprint(w, fmt.Sprintf(`[{"id":1,"name":"%s"}]`, fileName))
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/files/1",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`{"id":1,"name":"%s","contents":"%s"}`, fileName, fileContents))
		},
	)

	for _, data := range downloadTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdDownload(data.svcName, data.fileName, data.output, data.force, New(settings), services.New(settings))

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
	os.Remove("output.txt")
}

func TestDownloadForce(t *testing.T) {
	os.Remove("output.txt")
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
			fmt.Fprint(w, fmt.Sprintf(`[{"id":1,"name":"%s"}]`, fileName))
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/files/1",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`{"id":1,"name":"%s","contents":"%s"}`, fileName, fileContents))
		},
	)

	err := CmdDownload(test.SvcLabel, fileName, "output.txt", false, New(settings), services.New(settings))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	err = CmdDownload(test.SvcLabel, fileName, "output.txt", false, New(settings), services.New(settings))
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	err = CmdDownload(test.SvcLabel, fileName, "output.txt", true, New(settings), services.New(settings))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	os.Remove("output.txt")
}
