package invites

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

func CmdRm(inviteID string, ii IInvites) error {
	err := ii.Rm(inviteID)
	if err != nil {
		return err
	}
	logrus.Printf("Invite %s removed", inviteID)
	return nil
}

// Rm deletes an invite sent to a user. This invite must not already be
// accepted.
func (i *SInvites) Rm(inviteID string) error {
	headers := i.Settings.HTTPManager.GetHeaders(i.Settings.SessionToken, i.Settings.Version, i.Settings.Pod, i.Settings.UsersID)
	resp, statusCode, err := i.Settings.HTTPManager.Delete(nil, fmt.Sprintf("%s%s/orgs/%s/invites/%s", i.Settings.AuthHost, i.Settings.AuthHostVersion, i.Settings.OrgID, inviteID), headers)
	if err != nil {
		return err
	}
	return i.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
