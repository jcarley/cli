package worker

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
)

func CmdWorker(target, svcID string, iw IWorker) error {
	fmt.Printf("Initiating a background worker for Service: %s (procfile target = \"%s\")\n", svcID, target)
	err := iw.Start(target)
	if err != nil {
		return err
	}
	fmt.Println("Worker started.")
	return nil
}

// Start starts a Procfile target as a worker. Worker containers are intended
// to be short-lived, one-off tasks.
func (w *SWorker) Start(target string) error {
	helpers.InitiateWorker(target, w.Settings)
	return nil
}
