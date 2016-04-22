package supportids

import "github.com/Sirupsen/logrus"

func CmdSupportIDs(is ISupportIDs) error {
	envID, orgID, usersID, svcID, podID, err := is.SupportIDs()
	if err != nil {
		return err
	}
	logrus.Printf(`EnvironmentID:  %s
OrganizationID: %s
UsersID:        %s
ServiceID:      %s
PodID:          %s`, envID, orgID, usersID, svcID, podID)
	return nil
}

// SupportIDs prints out various IDs related to the associated environment to be
// used when contacting Catalyze support at support@catalyze.io.
func (s *SSupportIDs) SupportIDs() (string, string, string, string, string, error) {
	return s.Settings.EnvironmentID, s.Settings.OrgID, s.Settings.UsersID, s.Settings.ServiceID, s.Settings.Pod, nil
}
