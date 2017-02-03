package worker

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/lib/jobs"
)

func CmdDeploy(svcName, target string, iw IWorker, is services.IServices, ij jobs.IJobs) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"datica services list\" command.", svcName)
	}
	logrus.Printf("Initiating a worker for service %s (procfile target = \"%s\")", svcName, target)
	workers, err := iw.Retrieve(service.ID)
	if err != nil {
		return err
	}
	if _, ok := workers.Workers[target]; ok {
		logrus.Printf("Worker with target %s for service %s is already running, deploying another worker", target, svcName)
	}
	workers.Workers[target]++
	err = iw.Update(service.ID, workers)
	if err != nil {
		return err
	}
	err = ij.DeployTarget(target, service.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Successfully deployed a worker for service %s with target %s", svcName, target)
	return nil
}
