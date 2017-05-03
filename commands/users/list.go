package users

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/invites"
	"github.com/daticahealth/cli/models"
	"github.com/olekukonko/tablewriter"
)

func CmdList(myUsersID string, iu IUsers, ii invites.IInvites) error {
	orgGroups, err := ii.ListOrgGroups()
	if err != nil {
		return err
	}
	if orgGroups == nil || len(*orgGroups) == 0 {
		logrus.Println("No users found")
		return nil
	}
	data := [][]string{{"EMAIL", "GROUP(S)"}}
	members := make(map[string][]string)
	for _, group := range *orgGroups {
		groupMembers := group.Members
		for _, member := range *groupMembers {
			if _, ok := members[member.Email]; ok {
				members[member.Email] = append(members[member.Email], group.Name)
			} else {
				members[member.Email] = []string{group.Name}
			}
		}
	}
	for k, v := range members {
		data = append(data, []string{k, strings.Join(v, ", ")})
	}
	table := tablewriter.NewWriter(logrus.StandardLogger().Out)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.AppendBulk(data)
	table.Render()
	return nil
}

func (u *SUsers) List() (*[]models.OrgUser, error) {
	headers := u.Settings.HTTPManager.GetHeaders(u.Settings.SessionToken, u.Settings.Version, u.Settings.Pod, u.Settings.UsersID)
	resp, statusCode, err := u.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/orgs/%s/users", u.Settings.AuthHost, u.Settings.AuthHostVersion, u.Settings.OrgID), headers)
	if err != nil {
		return nil, err
	}
	var users []models.OrgUser
	err = u.Settings.HTTPManager.ConvertResp(resp, statusCode, &users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}
