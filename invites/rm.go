package invites

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

// RmInvite deletes an invite sent to a user. This invite must not already be
// accepted.
func RmInvite(inviteID string, settings *models.Settings) {
	helpers.SignIn(settings)
	helpers.DeleteInvite(inviteID, settings)
	fmt.Printf("Invite %s removed\n", inviteID)
}
