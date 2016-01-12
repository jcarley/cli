package whoami

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

// WhoAmI prints out your user ID. This can be used for adding users to
// environments via `catalyze adduser`, removing users from an environment
// via `catalyze rmuser`, when contacting Catalyze Support, etc.
func WhoAmI(settings *models.Settings) {
	helpers.SignIn(settings)
	fmt.Printf("user ID = %s\n", settings.UsersID)
}
