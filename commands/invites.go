package commands

import (
	"fmt"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
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

// ListInvites lists all pending invites for a given environment. After an
// invite is accepted, you can manage the users access with the `users`
// commands.
func ListInvites(settings *models.Settings) {
	helpers.SignIn(settings)
	invites := helpers.ListEnvironmentInvites(settings)
	fmt.Printf("Pending invites for %s:\n", settings.EnvironmentName)
	for _, invite := range *invites {
		fmt.Printf("\t%s %s\n", invite.Email, invite.Code)
	}
}

// RmInvite deletes an invite sent to a user. This invite must not already be
// accepted.
func RmInvite(inviteID string, settings *models.Settings) {
	helpers.SignIn(settings)
	helpers.DeleteInvite(inviteID, settings)
	fmt.Printf("Invite %s removed\n", inviteID)
}
