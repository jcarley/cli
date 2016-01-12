package db

import (
	"fmt"
	"sort"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

// SortedJobs is a wrapper for Jobs array in order to sort them by CreatedAt
// for the ListBackups command
type SortedJobs []models.Job

func (jobs SortedJobs) Len() int {
	return len(jobs)
}

func (jobs SortedJobs) Swap(i, j int) {
	jobs[i], jobs[j] = jobs[j], jobs[i]
}

func (jobs SortedJobs) Less(i, j int) bool {
	return jobs[i].CreatedAt < jobs[j].CreatedAt
}

// List lists the created backups for the service sorted from oldest to newest
func (d *SDb) List() error {
	service := helpers.RetrieveServiceByLabel(d.DatabaseName, d.Settings)
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"\n", d.DatabaseName)
	}
	jobs := helpers.ListBackups(service.ID, d.Page, d.PageSize, d.Settings)
	sort.Sort(SortedJobs(*jobs))
	for _, job := range *jobs {
		fmt.Printf("%s %s (status = %s)\n", job.ID, job.CreatedAt, job.Status)
	}
	if len(*jobs) == d.PageSize && d.Page == 1 {
		fmt.Println("(for older backups, try with --page 2 or adjust --page-size)")
	}
	if len(*jobs) == 0 && d.Page == 1 {
		fmt.Println("No backups created yet for this service.")
	}
	return nil
}
