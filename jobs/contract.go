package jobs

import "github.com/catalyzeio/cli/models"

// IJobs
type IJobs interface {
	Retrieve(jobID string) (*models.Job, error)
	RetrieveFromTaskID(taskID string) (*models.Job, error)
	RetrieveByStatus(status string) (*map[string]models.Job, error)
	RetrieveByType(jobType string, page, pageSize int) (*map[string]models.Job, error)
}

// SJobs is a concrete implementation of IJobs
type SJobs struct {
	Settings *models.Settings
}

// New returns an instance of IJobs
func New(settings *models.Settings) IJobs {
	return &SJobs{
		Settings: settings,
	}
}
