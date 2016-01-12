package supportids

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

// SupportIds prints out various IDs related to the associated environment to be
// used when contacting Catalyze support at support@catalyze.io.
func SupportIds(settings *models.Settings) {
	helpers.SignIn(settings)
	fmt.Printf(`EnvironmentID:  %s
UsersID:        %s
ServiceID:      %s
`, settings.EnvironmentID, settings.UsersID, settings.ServiceID)
}
