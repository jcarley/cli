package images

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/lib/images"
	"github.com/daticahealth/cli/test"
)

func TestListImagesNoRegistry(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())

	mux.HandleFunc("/environments/"+test.EnvID,
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`{"id":"%s","name":"%s","namespace":"%s","organizationId":"%s","docker_registry_enabled": false}`, test.EnvID, test.EnvName, test.Namespace, test.OrgID))
		},
	)

	err := cmdImageList(test.EnvID, environments.New(settings), images.New(settings))

	if err == nil {
		t.Fatal("Expected an error when the environment wasn't enabled to have a registry")
	}
}

func TestListImages(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())

	mux.HandleFunc("/environments/"+test.EnvID,
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`{"id":"%s","name":"%s","namespace":"%s","organizationId":"%s","docker_registry_enabled": true}`, test.EnvID, test.EnvName, test.Namespace, test.OrgID))
		},
	)

	mux.HandleFunc(fmt.Sprintf("/environments/"+test.EnvID+"/images"),
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, `["pod1234/a", "pod1234/b", "pod1234/c"]`)
		},
	)

	err := cmdImageList(test.EnvID, environments.New(settings), images.New(settings))

	if err != nil {
		t.Fatalf("Unexpected error while listing images: %v", err)
	}
}
