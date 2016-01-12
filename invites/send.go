package invites

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

// InviteUser invites a user by email to the associated environment. They do
// not need a Dashboard account prior to inviting them, but they must have a
// Dashboard account in order to accept the invitation.
func InviteUser(email string, settings *models.Settings) {
	helpers.YesNoPrompt(fmt.Sprintf("Are you sure you want to invite %s to your %s environment? (y/n) ", email, settings.EnvironmentName))
	helpers.SignIn(settings)
	helpers.CreateInvite(email, settings)
	fmt.Printf("%s has been invited!\n", email)
}
