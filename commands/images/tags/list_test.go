package tags

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/lib/images"
	"github.com/daticahealth/cli/test"
)

var listTagsTests = []struct {
	imageName string
	expectErr bool
}{
	{"hello", false},
	{"hello:tag", false},
	{fmt.Sprintf("%s/hello", test.Namespace), false},
	{fmt.Sprintf("%s/%s/hello", test.Registry, test.Namespace), false},
	{fmt.Sprintf("%s/invalid/hello", test.Registry), true},
}

func TestListTags(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID,
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`{"id":"%s","name":"%s","namespace":"%s","organizationId":"%s"}`, test.EnvID, test.EnvName, test.Namespace, test.OrgID))
		},
	)
	mux.HandleFunc(fmt.Sprintf("/environments/"+test.EnvID+"/images/"+test.Namespace+"%%2F"+"hello/tags"),
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, `["v1", "v2", "latest"]`)
		},
	)

	for _, data := range listTagsTests {
		err := cmdTagList(images.New(settings), environments.New(settings), settings.EnvironmentID, data.imageName)
		if (err != nil) != data.expectErr {
			t.Fatalf("Unexpected error while listing tags: %v", err)
		}
	}
}
