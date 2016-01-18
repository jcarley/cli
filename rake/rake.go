package rake

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/httpclient"
)

func CmdRake(taskName string, ir IRake) error {
	logrus.Printf("Executing Rake task: %s", taskName)
	err := ir.Run(taskName)
	if err != nil {
		return err
	}
	logrus.Println("Rake task output viewable in your logging server")
	return nil
}

// Run executes a rake task. This is only applicable for ruby-based
// applications.
func (r *SRake) Run(taskName string) error {
	rakeTask := map[string]string{
		"rake_command": taskName,
	}
	b, err := json.Marshal(rakeTask)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(r.Settings.SessionToken, r.Settings.Version, r.Settings.Pod)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/services/%s/brrgc/rake", r.Settings.PaasHost, r.Settings.PaasHostVersion, r.Settings.ServiceID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
