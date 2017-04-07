package deploykeys

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/test"
)

const (
	pubKeyPath = "cli_test.pub"
	keyName    = "my-key"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if err := createKeyFiles(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	statusCode := m.Run()
	cleanupKeyFiles()
	os.Exit(statusCode)
}

var addTests = []struct {
	name      string
	keyPath   string
	svcName   string
	expectErr bool
}{
	{keyName, pubKeyPath, test.SvcLabel, false},
	{"bad/key?name%", pubKeyPath, test.SvcLabel, true},
	{keyName, "~/invalid", test.SvcLabel, true},
	{keyName, pubKeyPath, "invalid-svc", true},
}

func TestAdd(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","label":"%s","type":"code"}]`, test.SvcID, test.SvcLabel))
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/ssh_keys",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "POST")
			w.WriteHeader(204)
		},
	)

	for _, data := range addTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdAdd(data.name, data.keyPath, data.svcName, New(settings), services.New(settings))

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}

func createKeyFiles() error {
	b, err := exec.Command("ssh-keygen", "-t", "rsa", "-b", "2048", "-N", "", "-f", strings.Replace(pubKeyPath, ".pub", "", 1)).CombinedOutput()
	fmt.Println(string(b))
	return err
}

func cleanupKeyFiles() error {
	err := os.Remove(pubKeyPath)
	if err == nil {
		err = os.Remove(strings.Replace(pubKeyPath, ".pub", "", 1))
	}
	return err
}
