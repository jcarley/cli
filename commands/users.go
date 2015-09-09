package commands

import (
	"fmt"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
)

// AddUser grants a user access to the associated environment. The ID of the
// user is required which can be found via `catalyze whoami`.
func AddUser(usersID string, settings *models.Settings) {
	helpers.SignIn(settings)
	helpers.AddUserToEnvironment(usersID, settings)
	fmt.Println("Added.")
}

// RmUser revokes a user's access to the associated environment. The ID of the
// user is required which can be found via `catalyze whoami`.
func RmUser(usersID string, settings *models.Settings) {
	helpers.SignIn(settings)
	helpers.RemoveUserFromEnvironment(usersID, settings)
	fmt.Println("Removed.")
}

// ListUsers lists all users who have access to the associated environment.
func ListUsers(settings *models.Settings) {
	helpers.SignIn(settings)
	envUsers := helpers.ListEnvironmentUsers(settings)
	for _, userID := range envUsers.Users {
		if userID == settings.UsersID {
			fmt.Printf("%s (you)\n", userID)
		} else {
			defer fmt.Printf("%s\n", userID)
		}
	}
}

// WhoAmI prints out your user ID. This can be used for adding users to
// environments via `catalyze adduser`, removing users from an environment
// via `catalyze rmuser`, when contacting Catalyze Support, etc.
func WhoAmI(settings *models.Settings) {
	helpers.SignIn(settings)
	fmt.Printf("user ID = %s\n", settings.UsersID)
}
