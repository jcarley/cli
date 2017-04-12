package files

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/daticahealth/cli/test"
)

const (
	filePath = "test_file.txt"
	fileName = "file.txt"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if err := createFiles(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	statusCode := m.Run()
	cleanupFiles()
	os.Exit(statusCode)
}

var createTests = []struct {
	svcID     string
	filePath  string
	name      string
	mode      string
	expectErr bool
}{
	{test.SvcID, filePath, fileName, "0644", false},
	{"invalid-svc", filePath, fileName, "0644", true},
	{test.SvcID, "invalid-file-path", fileName, "0644", true},
}

func TestCreate(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/files",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "POST")
			fmt.Fprint(w, `{"id":1,"contents":""}`)
		},
	)

	ifiles := New(settings)

	for _, data := range createTests {
		t.Logf("Data: %+v", data)

		// test
		_, err := ifiles.Create(data.svcID, data.filePath, data.name, data.mode)

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}

func createFiles() error {
	return ioutil.WriteFile(filePath, []byte("text"), 0644)
}

func cleanupFiles() error {
	return os.Remove(filePath)
}
