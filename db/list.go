package db

import (
	"fmt"
	"os"
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

// ListBackups lists the created backups for the service sorted from oldest to newest
func ListBackups(serviceLabel string, page int, pageSize int, settings *models.Settings) {
	helpers.SignIn(settings)
	service := helpers.RetrieveServiceByLabel(serviceLabel, settings)
	if service == nil {
		fmt.Printf("Could not find a service with the label \"%s\"\n", serviceLabel)
		os.Exit(1)
	}
	jobs := helpers.ListBackups(service.ID, page, pageSize, settings)
	sort.Sort(SortedJobs(*jobs))
	for _, job := range *jobs {
		fmt.Printf("%s %s (status = %s)\n", job.ID, job.CreatedAt, job.Status)
	}
	if len(*jobs) == pageSize && page == 1 {
		fmt.Println("(for older backups, try with --page 2 or adjust --page-size)")
	}
	if len(*jobs) == 0 && page == 1 {
		fmt.Println("No backups created yet for this service.")
	}
}
