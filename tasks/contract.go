package tasks

import "github.com/catalyzeio/cli/models"

// ITasks
type ITasks interface {
	PollForStatus(task *models.Task) (string, error)
	PollForConsole(task *models.Task, service *models.Service) (string, error)
}

// STasks is a concrete implementation of ITasks
type STasks struct {
	Settings *models.Settings
}

// New returns an instance of ITasks
func New(settings *models.Settings) ITasks {
	return &STasks{
		Settings: settings,
	}
}
