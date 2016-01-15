package invites

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/prompts"
)

func CmdSend(email, envName string, ii IInvites, ip prompts.IPrompts) error {
	err := ip.YesNo(fmt.Sprintf("Are you sure you want to invite %s to your %s environment? (y/n) ", email, envName))
	if err != nil {
		return err
	}
	return ii.Send(email)
}

// Send invites a user by email to the associated environment. They do
// not need a Dashboard account prior to inviting them, but they must have a
// Dashboard account in order to accept the invitation.
func (i *SInvites) Send(email string) error {
	helpers.CreateInvite(email, i.Settings)
	fmt.Printf("%s has been invited!\n", email)
	return nil
}
