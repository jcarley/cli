package invites

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/prompts"
)

func CmdSend(email, envName string, ii IInvites, ip prompts.IPrompts) error {
	err := ip.YesNo(fmt.Sprintf("Are you sure you want to invite %s to your %s environment? (y/n) ", email, envName))
	if err != nil {
		return err
	}
	err = ii.Send(email)
	if err != nil {
		return err
	}
	logrus.Printf("%s has been invited!", email)
	return nil
}

// Send invites a user by email to the associated environment. They do
// not need a Dashboard account prior to inviting them, but they must have a
// Dashboard account in order to accept the invitation.
func (i *SInvites) Send(email string) error {
	inv := models.Invite{
		Email: email,
	}
	b, err := json.Marshal(inv)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(i.Settings.SessionToken, i.Settings.Version, i.Settings.Pod)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/invites", i.Settings.PaasHost, i.Settings.PaasHostVersion, i.Settings.EnvironmentID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
