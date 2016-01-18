package invites

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
)

func CmdList(envName string, ii IInvites) error {
	invts, err := ii.List()
	if err != nil {
		return err
	}
	if len(*invts) == 0 {
		logrus.Printf("There are no pending invites for %s", envName)
		return nil
	}
	logrus.Printf("Pending invites for %s:", envName)
	for _, invite := range *invts {
		logrus.Printf("\t%s %s", invite.Email, invite.Code)
	}
	return nil
}

// List lists all pending invites for a given environment. After an
// invite is accepted, you can manage the users access with the `users`
// commands.
func (i *SInvites) List() (*[]models.Invite, error) {
	headers := httpclient.GetHeaders(i.Settings.SessionToken, i.Settings.Version, i.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/invites", i.Settings.PaasHost, i.Settings.PaasHostVersion, i.Settings.EnvironmentID), headers)
	if err != nil {
		return nil, err
	}
	var invites []models.Invite
	err = httpclient.ConvertResp(resp, statusCode, &invites)
	if err != nil {
		return nil, err
	}
	return &invites, nil
}
