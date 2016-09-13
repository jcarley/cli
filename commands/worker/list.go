package worker

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/models"
	"github.com/olekukonko/tablewriter"
)

func CmdList(svcName string, iw IWorker, is services.IServices, ij jobs.IJobs) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	workers, err := iw.Retrieve(service.ID)
	if err != nil {
		return err
	}
	jobs, err := ij.RetrieveByType(service.ID, "worker", 1, 1000)
	if err != nil {
		return err
	}
	type workerJob struct {
		scale   int
		running int
	}
	var workerJobs = map[string]*workerJob{}
	for target, scale := range workers.Workers {
		workerJobs[target] = &workerJob{scale, 0}
	}
	for _, j := range *jobs {
		if _, ok := workerJobs[j.Target]; !ok {
			workerJobs[j.Target] = &workerJob{0, 0}
		}
		workerJobs[j.Target].running = 1
	}

	if len(workerJobs) == 0 {
		logrus.Printf("No workers found for service %s", svcName)
		return nil
	}

	data := [][]string{{"TARGET", "SCALE", "RUNNING JOBS"}}
	for target, wj := range workerJobs {
		data = append(data, []string{target, fmt.Sprintf("%d", wj.scale), fmt.Sprintf("%d", wj.running)})
	}

	table := tablewriter.NewWriter(logrus.StandardLogger().Out)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.AppendBulk(data)
	table.Render()
	return nil
}

func (w *SWorker) Retrieve(svcID string) (*models.Workers, error) {
	return nil, nil
}
