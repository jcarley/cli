package commands

import (
	"fmt"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
)

// Rake executes a rake task. This is only applicable for ruby-based
// applications.
func Rake(taskName string, settings *models.Settings) {
	helpers.SignIn(settings)
	fmt.Printf("Executing Rake task: %s\n", taskName)
	helpers.InitiateRakeTask(taskName, settings)
	fmt.Println("Rake task output viewable in your logging server")
}
