package invites

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

// ListInvites lists all pending invites for a given environment. After an
// invite is accepted, you can manage the users access with the `users`
// commands.
func ListInvites(settings *models.Settings) {
	helpers.SignIn(settings)
	invites := helpers.ListEnvironmentInvites(settings)
	if len(*invites) == 0 {
		fmt.Printf("There are no pending invites for %s\n", settings.EnvironmentName)
		return
	}
	fmt.Printf("Pending invites for %s:\n", settings.EnvironmentName)
	for _, invite := range *invites {
		fmt.Printf("\t%s %s\n", invite.Email, invite.Code)
	}
}
