package worker

import (
	"fmt"

	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/lib/jobs"
)

func CmdWorker(svcName, defaultSvcID, target string, iw IWorker, is services.IServices, ij jobs.IJobs) error {
	if svcName != "" {
		service, err := is.RetrieveByLabel(svcName)
		if err != nil {
			return err
		}
		if service == nil {
			return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"datica services list\" command.", svcName)
		}
		svcName = service.Label
	}
	return CmdDeploy(svcName, target, iw, is, ij)
}
