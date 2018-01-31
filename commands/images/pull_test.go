package images

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/models"
	"github.com/daticahealth/cli/test"
)

func pullSetup(t *testing.T, mux *http.ServeMux, baseURL string) (string, string) {
	notary := baseURL
	registry := strings.TrimPrefix(baseURL, "http://")
	test.AddRegistry("default", registry)
	test.AddNotary("default", notary)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		test.AssertEquals(t, r.Method, "GET")
	})

	mux.HandleFunc("/environments/"+test.EnvID,
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`{"id":"%s","name":"%s","namespace":"%s","organizationId":"%s","docker_registry_enabled": true}`, test.EnvID, test.EnvName, test.Namespace, test.OrgID))
		},
	)

	return registry, notary
}

func TestPullImage(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	registry, notary := pullSetup(t, mux, baseURL.String())
	user := models.User{}
	fakeImages := &test.FakeImages{Settings: settings}
	imageTag := fmt.Sprintf("%s:%s", test.Image, test.Tag)
	fullImageName := strings.Join([]string{registry, test.Namespace, imageTag}, "/")
	test.SetRemoteImages([]string{fullImageName})

	err := cmdImagePull(test.EnvID, imageTag, &user, environments.New(settings), fakeImages)

	cleanupErr := test.DeleteLocalRepo(test.Namespace, test.Image, registry, notary)

	if err != nil {
		t.Fatalf("Unexpected error while pushing image: %v", err)
	} else if cleanupErr != nil {
		t.Errorf("Unexpected error while cleaning up trust repo (test succeeded otherwise): %v", cleanupErr)
	}
}

func TestPullBadImage(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	pullSetup(t, mux, baseURL.String())
	user := models.User{}
	fakeImages := &test.FakeImages{Settings: settings}
	imageTag := fmt.Sprintf("%s:%s", test.Image, test.Tag)

	err := cmdImagePull(test.EnvID, imageTag, &user, environments.New(settings), fakeImages)

	if err == nil {
		t.Fatalf("Expected error: %v", test.ImageDoesNotExist)
	} else if err.Error() != test.ImageDoesNotExist {
		t.Fatalf("Expected error: %v\nGot error: %v", test.ImageDoesNotExist, err)
	}
}
