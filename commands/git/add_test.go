package git

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/test"
)

var addTests = []struct {
	svcLabel    string
	remote      string
	removeFirst bool
	force       bool
	expectErr   bool
}{
	{test.SvcLabel, "datica", true, false, false},
	{"invalid-svc", "datica", true, false, true},
	{test.SvcLabel, "datica", false, false, true},
	{test.SvcLabel, "datica", false, true, false},
}

func TestAdd(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprintf(w, `[{"id":"%s","label":"%s","source":"git@github.com/github/github.git"}]`, test.SvcID, test.SvcLabel)
		},
	)

	for _, data := range addTests {
		t.Logf("Data: %+v", data)

		// test
		ig := New()
		if data.removeFirst {
			ig.Rm(data.remote)
		}
		err := CmdAdd(data.svcLabel, data.remote, data.force, ig, services.New(settings))

		// assert
		if err != nil {
			if !data.expectErr {
				t.Errorf("Unexpected error: %s", err)
			}
			continue
		}

		remotes, err := ig.List()
		if err != nil {
			t.Errorf("Failed to list remotes: %s", err)
			continue
		}
		found := false
		for _, r := range remotes {
			if r == data.remote {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Git remote was not added. Expected %s but found %v", data.remote, remotes)
		}
	}
}

func TestAddNoGitRepo(t *testing.T) {
	settings := test.GetSettings("")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %s", err)
	}
	t.Log(wd)
	defer os.Chdir(wd)

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Failed to make temp directory: %s", err)
	}

	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("Failed to change working directory: %s", err)
	}

	err = CmdAdd(test.SvcLabel, "datica", false, New(), services.New(settings))

	// assert
	if err == nil {
		t.Fatalf("Expected error adding a remote without a git repo")
	}
}
