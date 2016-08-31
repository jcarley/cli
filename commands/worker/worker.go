package worker

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
)

func CmdWorker(svcName, defaultSvcID, target string, iw IWorker, is services.IServices) error {
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
	logrus.Printf("Initiating a background worker for Service: %s (procfile target = \"%s\")", defaultSvcID, target)
	err := iw.Start(defaultSvcID, target)
	if err != nil {
		return err
	}
	logrus.Println("Worker started.")
	return nil
}

// Start starts a Procfile target as a worker. Worker containers are intended
// to be short-lived, one-off tasks.
func (w *SWorker) Start(svcID, target string) error {
	worker := map[string]string{
		"target": target,
	}
	b, err := json.Marshal(worker)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(w.Settings.SessionToken, w.Settings.Version, w.Settings.Pod, w.Settings.UsersID)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/worker", w.Settings.PaasHost, w.Settings.PaasHostVersion, w.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
