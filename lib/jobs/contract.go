package jobs

import "github.com/catalyzeio/cli/models"

// IJobs
type IJobs interface {
	Delete(jobID, svcID string) error
	Deploy(redeploy bool, releaseName, target, svcID string) error
	DeployRelease(releaseName, svcID string) error
	DeployTarget(target, svcID string) error
	Redeploy(svcID string) error
	Retrieve(jobID, svcID string, includeSpec bool) (*models.Job, error)
	RetrieveByStatus(svcID, status string) (*[]models.Job, error)
	RetrieveByType(svcID, jobType string, page, pageSize int) (*[]models.Job, error)
	RetrieveByTarget(svcID, target string, page, pageSize int) (*[]models.Job, error)
	PollForStatus(statuses []string, jobID, svcID string) (string, error)
	PollTillFinished(jobID, svcID string) (string, error)
	List(svcID string, page, pageSize int) (*[]models.Job, error)
	WaitToAppear(jobID, svcID string) error
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
