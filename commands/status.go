package commands

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
	"github.com/pmylund/sortutil"
)

// Status prints out an environment healthcheck. The status of the environment
// and every service in the environment is printed out.
func Status(settings *models.Settings) {
	helpers.SignIn(settings)
	env := helpers.RetrieveEnvironment("pod", settings)
	fmt.Printf("%s (environment ID = %s):\n", env.Data.Name, env.ID)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, '\t', 0)
	fmt.Fprintln(w, "ID\tType\tStatus\tCreated At")
	r := *env.Data.Services
	sortutil.AscByField(r, "Label")
	var latestBuild models.Job
	for _, service := range r {
		if service.Type != "utility" && service.Type != "" {
			jobs := helpers.RetrieveAllJobs(service.ID, settings)
			for jobID, job := range *jobs {
				if job.Status == "running" {
					const dateForm = "2006-01-02T15:04:05"
					t, _ := time.Parse(dateForm, job.CreatedAt)

					displayType := service.Label
					if job.Type != "deploy" {
						displayType = fmt.Sprintf("%s (%s)", displayType, job.Type)
					}
					fmt.Fprintln(w, jobID[:8]+"\t"+displayType+"\t"+job.Status+"\t"+t.Local().Format(time.Stamp))
				}
				if job.Status == "finished" && job.Type == "build" && job.CreatedAt > latestBuild.CreatedAt {
					latestBuild = job
					latestBuild.Type = job.Type
					latestBuild.Status = job.Status
					latestBuild.ID = jobID
				}
			}
			if service.Type == "code" {
				const dateForm = "2006-01-02T15:04:05"
				t, _ := time.Parse(dateForm, latestBuild.CreatedAt)
				displayType := service.Label
				displayType = fmt.Sprintf("%s (%s)", displayType, latestBuild.Type)
				fmt.Fprintln(w, latestBuild.ID[:8]+"\t"+displayType+"\t"+latestBuild.Status+"\t"+t.Local().Format(time.Stamp))
			}
		}
	}
	w.Flush()
}
