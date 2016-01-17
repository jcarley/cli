package invites

import (
	"fmt"

	"github.com/catalyzeio/cli/httpclient"
)

func CmdRm(inviteID string, ii IInvites) error {
	err := ii.Rm(inviteID)
	if err != nil {
		return err
	}
	fmt.Printf("Invite %s removed\n", inviteID)
	return nil
}

// Rm deletes an invite sent to a user. This invite must not already be
// accepted.
func (i *SInvites) Rm(inviteID string) error {
	headers := httpclient.GetHeaders(i.Settings.APIKey, i.Settings.SessionToken, i.Settings.Version, i.Settings.Pod)
	resp, statusCode, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/environments/%s/invites/%s", i.Settings.PaasHost, i.Settings.PaasHostVersion, i.Settings.EnvironmentID, inviteID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
