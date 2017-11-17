package tags

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/daticahealth/cli/lib/images"
	"github.com/daticahealth/cli/test"
)

func TestListTags(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())

	mux.HandleFunc(fmt.Sprintf("/environments/"+test.EnvID+"/images/test1234%%2Fhello/tags"),
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, `["v1", "v2", "latest"]`)
		},
	)

	err := cmdTagList(images.New(settings), "test1234/hello")

	if err != nil {
		t.Fatalf("Unexpected error while listing tags: %v", err)
	}
}
