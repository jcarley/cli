package supportids

import "github.com/Sirupsen/logrus"

func CmdSupportIDs(is ISupportIDs) error {
	envID, usersID, svcID, err := is.SupportIDs()
	if err != nil {
		return err
	}
	logrus.Printf(`EnvironmentID:  %s
UsersID:        %s
ServiceID:      %s`, envID, usersID, svcID)
	return nil
}

// SupportIDs prints out various IDs related to the associated environment to be
// used when contacting Catalyze support at support@catalyze.io.
func (s *SSupportIDs) SupportIDs() (string, string, string, error) {
	return s.Settings.EnvironmentID, s.Settings.UsersID, s.Settings.ServiceID, nil
}
