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
	orgUsers, err := iu.List()
	if err != nil {
		return err
	}
	if orgUsers == nil || len(*orgUsers) == 0 {
		logrus.Println("No users found")
		return nil
	}
	orgGroups, err := ii.ListOrgGroups()
	if err != nil {
		return err
	}
	members := make(map[string][]string)
	for _, group := range *orgGroups {
		groupMembers := group.Members
		for _, member := range *groupMembers {
			members[member.Email] = append(members[member.Email], group.Name)
		}
	}
	data := [][]string{{"EMAIL", "GROUP(S)"}}
	for _, user := range *orgUsers {
		if val, ok := members[user.Email]; ok {
			data = append(data, []string{user.Email, strings.Join(val, ", ")})
		} else {
			data = append(data, []string{user.Email, "none"})
		}
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
