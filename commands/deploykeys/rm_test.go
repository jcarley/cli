package deploykeys

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/test"
)

const sampleKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCykOkMAXDMKuz/6x1bVT/4Cz6JjDrEnkjbwvFq6Gp9NLW79vjcHJdhkeYEhnhtmNf62PVP2lHgmkxuk0OFl5mZg+SZeNep/cSfKdp99KjG2cGWd2XDDxwQLG8JyUcLRZ+1Q653lncwc6vL+hmBCvQ4gQhx9OA+XzNk064BQb/BCMWyLvQXAXr2dGKs1jIDV/CsMBgXmz4KjnYuuBYM3o44MeYw9fMPFz5J+i/sZUDPdXUAGGQd8cYrpxfNI4Qn33839uspa0eCCT6iMuCBX9heJs6CQ77vl+5TFdZjr+xahegcnDGtHyxUL76Jwm22bJ3jFJCDWYr1/vTZvJy6+wf9"

var rmTests = []struct {
	keyName   string
	svcName   string
	expectErr bool
}{
	{keyName, test.SvcLabel, false},
	{"invalid-key", test.SvcLabel, true},
	{keyName, "invalid-svc", true},
}

func TestRm(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","label":"%s","type":"code"}]`, test.SvcID, test.SvcLabel))
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/ssh_keys/"+keyName+"/type/ssh",
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "DELETE")
			w.WriteHeader(204)
		},
	)

	for _, data := range rmTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdRm(data.keyName, data.svcName, New(settings), services.New(settings))

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}
