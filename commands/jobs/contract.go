package jobs

import "github.com/catalyzeio/cli/models"

// IJobs
type IJobs interface {
	Retrieve(jobID, svcID string) (*models.Job, error)
	RetrieveByStatus(status string) (*[]models.Job, error)
	RetrieveByType(jobType string, page, pageSize int) (*[]models.Job, error)
	PollForStatus(jobID, svcID string) (string, error)
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
