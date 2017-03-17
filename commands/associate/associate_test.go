package associate

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/commands/git"
	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/models"
	"github.com/daticahealth/cli/test"
)

var associateTests = []struct {
	envName   string
	svcName   string
	alias     string
	remote    string
	expectErr bool
}{
	{test.EnvName, test.SvcLabel, "", "datica", false},
	{test.EnvName, test.SvcLabel, "", "custom", false},
	{test.EnvName, test.SvcLabel, test.Alias, "datica", false},
	{test.EnvName, "bad-svc", "", "datica", true},
	{"bad-env", test.SvcLabel, "", "datica", true},
}

func TestAssociate(t *testing.T) {
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
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			test.AssertEquals(t, r.Header.Get("X-Pod-ID"), test.Pod)
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","type":"code","label":"%s","source":"ssh://git@datica.com"}]`, test.SvcID, test.SvcLabel))
		},
	)
	mux.HandleFunc("/environments/"+test.EnvIDAlt+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			test.AssertEquals(t, r.Header.Get("X-Pod-ID"), test.PodAlt)
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","type":"code","label":"%s","source":"ssh://git@datica.com"}]`, test.SvcIDAlt, test.SvcLabelAlt))
		},
	)

	for _, data := range associateTests {
		t.Logf("Data: %+v", data)

		// reset
		settings.Environments = map[string]models.AssociatedEnv{}

		// test
		err := CmdAssociate(data.envName, data.svcName, data.alias, data.remote, false, New(settings), git.New(), environments.New(settings), services.New(settings))

		// assertions
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
		expectedEnvs := map[string]models.AssociatedEnv{}
		if !data.expectErr {
			name := data.alias
			if name == "" {
				name = data.envName
			}
			expectedEnvs[name] = models.AssociatedEnv{
				Name:          test.EnvName,
				EnvironmentID: test.EnvID,
				ServiceID:     test.SvcID,
				Directory:     "",
				OrgID:         test.OrgID,
				Pod:           test.Pod,
			}
		}
		actual := map[string]models.AssociatedEnv{}
		for envName, env := range settings.Environments {
			env.Directory = ""
			actual[envName] = env
		}
		if !reflect.DeepEqual(expectedEnvs, actual) {
			t.Errorf("Associated environment not added to settings object correctly.\nExpected: %+v.\nFound:    %+v", expectedEnvs, settings.Environments)
		}
	}
}

func TestAssociateWithPodErrors(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	settings.Environments = map[string]models.AssociatedEnv{}

	mux.HandleFunc("/environments",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			if r.Header.Get("X-Pod-ID") == test.Pod {
				fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","name":"%s","namespace":"%s","organizationId":"%s"}]`, test.EnvID, test.EnvName, test.Namespace, test.OrgID))
			} else {
				http.Error(w, `{"title":"Error","description":"error","code":1}`, 400)
			}
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			test.AssertEquals(t, r.Header.Get("X-Pod-ID"), test.Pod)
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","type":"code","label":"%s","source":"ssh://git@datica.com"}]`, test.SvcID, test.SvcLabel))
		},
	)

	// test
	err := CmdAssociate(test.EnvName, test.SvcLabel, "", "datica", false, New(settings), git.New(), environments.New(settings), services.New(settings))

	// assert
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	expectedEnvs := map[string]models.AssociatedEnv{
		test.EnvName: models.AssociatedEnv{
			Name:          test.EnvName,
			EnvironmentID: test.EnvID,
			ServiceID:     test.SvcID,
			Directory:     "",
			OrgID:         test.OrgID,
			Pod:           test.Pod,
		},
	}
	actual := map[string]models.AssociatedEnv{}
	for envName, env := range settings.Environments {
		env.Directory = ""
		actual[envName] = env
	}
	if !reflect.DeepEqual(expectedEnvs, actual) {
		t.Errorf("Associated environment not added to settings object correctly.\nExpected: %+v.\nFound:    %+v", expectedEnvs, settings.Environments)
	}
}
