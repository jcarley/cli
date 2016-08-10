package rake

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
)

func CmdRake(svcName, taskName, defaultSvcID string, ir IRake, is services.IServices) error {
	if svcName != "" {
		service, err := is.RetrieveByLabel(svcName)
		if err != nil {
			return err
		}
		if service == nil {
			return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services\" command.", svcName)
		}
		defaultSvcID = service.ID
	}
	logrus.Printf("Executing Rake task: %s", taskName)
	err := ir.Run(taskName, defaultSvcID)
	if err != nil {
		return err
	}
	logrus.Println("Rake task output viewable in your logging dashboard")
	return nil
}

// Run executes a rake task. This is only applicable for ruby-based
// applications.
func (r *SRake) Run(taskName, svcID string) error {
	rakeTask := map[string]string{
		"command": taskName,
	}
	b, err := json.Marshal(rakeTask)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(r.Settings.SessionToken, r.Settings.Version, r.Settings.Pod, r.Settings.UsersID)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/rake", r.Settings.PaasHost, r.Settings.PaasHostVersion, r.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
