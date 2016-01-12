package logout

import (
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/models"
)

// Logout clears the stored user information from the local machine. This does
// not remove environment data.
func Logout(settings *models.Settings) {
	settings.SessionToken = ""
	settings.UsersID = ""
	config.SaveSettings(settings)
}
