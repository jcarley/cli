package rake

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

// Rake executes a rake task. This is only applicable for ruby-based
// applications.
func Rake(taskName string, settings *models.Settings) {
	helpers.SignIn(settings)
	fmt.Printf("Executing Rake task: %s\n", taskName)
	helpers.InitiateRakeTask(taskName, settings)
	fmt.Println("Rake task output viewable in your logging server")
}
