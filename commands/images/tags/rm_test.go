package tags

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/lib/images"
	"github.com/daticahealth/cli/lib/prompts"
	"github.com/daticahealth/cli/test"
)

type tPrompts struct {
	*prompts.SPrompts
	accept bool
}

func New(accept bool) *tPrompts {
	tp := &tPrompts{}
	tp.accept = accept
	return tp
}

func (tp *tPrompts) YesNo(msg, prompt string) error {
	if tp.accept {
		return nil
	}
	return errors.New("declined")
}

func setupMux(mux *http.ServeMux, t *testing.T, successfulDelete bool) {
	mux.HandleFunc("/environments/"+test.EnvID,
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`{"id":"%s","name":"%s","namespace":"%s","organizationId":"%s"}`, test.EnvID, test.EnvName, test.Namespace, test.OrgID))
		},
	)
	mux.HandleFunc(fmt.Sprintf("/environments/"+test.EnvID+"/images/"+test.Namespace+"%%2F"+test.Image+"/tags"),
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, `["v1", "v2", "latest"]`)
		},
	)
	if successfulDelete {
		mux.HandleFunc(fmt.Sprintf("/environments/"+test.EnvID+"/images/"+test.Namespace+"%%2F"+test.Image+"/tags/"+test.Tag),
			func(w http.ResponseWriter, r *http.Request) {
				test.AssertEquals(t, r.Method, "DELETE")
				w.WriteHeader(204)
			},
		)
	}
}

var deleteTagTests = []struct {
	imageName string
	expectErr bool
}{
	// SUCCEED
	{fmt.Sprintf("%s:%s", test.Image, test.Tag), false},
	{fmt.Sprintf("%s/%s:%s", test.Namespace, test.Image, test.Tag), false},
	{fmt.Sprintf("%s/%s/%s:%s", test.Registry, test.Namespace, test.Image, test.Tag), false},
	// FAIL
	{fmt.Sprintf("%s", test.Image), true},
	{fmt.Sprintf("%s/%s", test.Namespace, test.Image), true},
	{fmt.Sprintf("%s/%s/%s", test.Registry, test.Namespace, test.Image), true},
	{fmt.Sprintf("invalid/%s", test.Image), true},
	{fmt.Sprintf("%s/invalid/%s", test.Registry, test.Image), true},
}

func TestDeleteTag(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	setupMux(mux, t, true)

	for _, data := range deleteTagTests {
		err := cmdTagDelete(images.New(settings), New(true), environments.New(settings), settings.EnvironmentID, data.imageName)

		if (err != nil) != data.expectErr {
			t.Fatalf("Unexpected error while deleting tag: %v", err)
		}
	}
}

func TestDeleteDecline(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	setupMux(mux, t, true)

	err := cmdTagDelete(images.New(settings), New(false), environments.New(settings), settings.EnvironmentID, fmt.Sprintf("%s/%s:%s", test.Namespace, test.Image, test.Tag))

	if err == nil {
		t.Fatalf("Expected error while deleting tag: %v", err)
	}
}

func TestDeleteTagNotInList(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	setupMux(mux, t, false)

	err := cmdTagDelete(images.New(settings), New(false), environments.New(settings), settings.EnvironmentID, fmt.Sprintf("%s/%s:%s", test.Namespace, test.Image, "v3"))

	if err == nil {
		t.Fatalf("Expected error while deleting tag: %v", err)
	}
}

func TestDeleteTagIsRelease(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	setupMux(mux, t, false)

	mux.HandleFunc(fmt.Sprintf("/environments/"+test.EnvID+"/images/"+test.Namespace+"%%2F"+test.Image+"/tags/"+test.Tag),
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "DELETE")
			w.WriteHeader(400)
			fmt.Fprint(w, `{"code": 98005, "title": "Tag is release", "description": "The tag is a release"}`)
		},
	)

	err := cmdTagDelete(images.New(settings), New(false), environments.New(settings), settings.EnvironmentID, fmt.Sprintf("%s/%s:%s", test.Namespace, test.Image, test.Tag))

	if err == nil {
		t.Fatalf("Expected error while deleting tag: %v", err)
	}
}
