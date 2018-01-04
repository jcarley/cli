package tags

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

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

func TestDeleteTag(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())

	mux.HandleFunc(fmt.Sprintf("/environments/"+test.EnvID+"/images/test1234%%2Fhello/tags"),
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, `["v1", "v2", "latest"]`)
		},
	)

	mux.HandleFunc(fmt.Sprintf("/environments/"+test.EnvID+"/images/test1234%%2Fhello/tags/v1"),
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "DELETE")
			w.WriteHeader(204)
		},
	)

	err := cmdTagDelete(images.New(settings), New(true), "test1234/hello", "v1")

	if err != nil {
		t.Fatalf("Unexpected error while deleting tag: %v", err)
	}
}

func TestDeleteDecline(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())

	mux.HandleFunc(fmt.Sprintf("/environments/"+test.EnvID+"/images/test1234%%2Fhello/tags"),
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, `["v1", "v2", "latest"]`)
		},
	)

	mux.HandleFunc(fmt.Sprintf("/environments/"+test.EnvID+"/images/test1234%%2Fhello/tags/v1"),
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "DELETE")
			w.WriteHeader(204)
		},
	)

	err := cmdTagDelete(images.New(settings), New(false), "test1234/hello", "v1")

	if err == nil {
		t.Fatalf("Expected error while deleting tag: %v", err)
	}
}

func TestDeleteTagNotInList(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())

	mux.HandleFunc(fmt.Sprintf("/environments/"+test.EnvID+"/images/test1234%%2Fhello/tags"),
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, `["v1", "v2", "latest"]`)
		},
	)

	err := cmdTagDelete(images.New(settings), New(false), "test1234/hello", "v3")

	if err == nil {
		t.Fatalf("Expected error while deleting tag: %v", err)
	}
}

func TestDeleteTagIsRelease(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())

	mux.HandleFunc(fmt.Sprintf("/environments/"+test.EnvID+"/images/test1234%%2Fhello/tags"),
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, `["v1", "v2", "latest"]`)
		},
	)

	mux.HandleFunc(fmt.Sprintf("/environments/"+test.EnvID+"/images/test1234%%2Fhello/tags/v1"),
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "DELETE")
			w.WriteHeader(400)
			fmt.Fprint(w, `{"code": 98005, "title": "Tag is release", "description": "The tag is a release"}`)
		},
	)

	err := cmdTagDelete(images.New(settings), New(false), "test1234/hello", "v1")

	if err == nil {
		t.Fatalf("Expected error while deleting tag: %v", err)
	}
}
