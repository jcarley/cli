package invites

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

func CmdList(envName string, ii IInvites) error {
	invts, err := ii.List()
	if err != nil {
		return err
	}
	if len(*invts) == 0 {
		fmt.Printf("There are no pending invites for %s\n", envName)
		return nil
	}
	fmt.Printf("Pending invites for %s:\n", envName)
	for _, invite := range *invts {
		fmt.Printf("\t%s %s\n", invite.Email, invite.Code)
	}
	return nil
}

// List lists all pending invites for a given environment. After an
// invite is accepted, you can manage the users access with the `users`
// commands.
func (i *SInvites) List() (*[]models.Invite, error) {
	return helpers.ListEnvironmentInvites(i.Settings), nil
}
