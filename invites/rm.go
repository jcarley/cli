package invites

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
)

func CmdRm(inviteID string, ii IInvites) error {
	return ii.Rm(inviteID)
}

// Rm deletes an invite sent to a user. This invite must not already be
// accepted.
func (i *SInvites) Rm(inviteID string) error {
	helpers.DeleteInvite(inviteID, i.Settings)
	fmt.Printf("Invite %s removed\n", inviteID)
	return nil
}
