package rake

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/catalyzeio/cli/httpclient"
)

func CmdRake(taskName string, ir IRake) error {
	fmt.Printf("Executing Rake task: %s\n", taskName)
	err := ir.Run(taskName)
	if err != nil {
		return err
	}
	fmt.Println("Rake task output viewable in your logging server")
	return nil
}

// Rake executes a rake task. This is only applicable for ruby-based
// applications.
func (r *SRake) Run(taskName string) error {
	rakeTask := map[string]string{}
	b, err := json.Marshal(rakeTask)
	if err != nil {
		return err
	}
	encodedTaskName, err := url.Parse(taskName)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(r.Settings.APIKey, r.Settings.SessionToken, r.Settings.Version, r.Settings.Pod)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/rake/%s", r.Settings.PaasHost, r.Settings.PaasHostVersion, r.Settings.EnvironmentID, r.Settings.ServiceID, encodedTaskName), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
