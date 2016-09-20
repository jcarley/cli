package redeploy

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/jobs"
)

func CmdRedeploy(svcName string, ij jobs.IJobs, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services\" command.", svcName)
	}
	logrus.Printf("Redeploying %s (ID = %s)", svcName, service.ID)
	err = ij.Redeploy(service.ID)
	if err != nil {
		return err
	}
	logrus.Println("Redeploy successful! Check the status with \"catalyze status\" and your logging dashboard for updates")
	return nil
}
