package whoami

import "github.com/Sirupsen/logrus"

func CmdWhoAmI(w IWhoAmI) error {
	usersID, err := w.WhoAmI()
	if err != nil {
		return err
	}
	logrus.Printf("user ID = %s", usersID)
	return nil
}

// WhoAmI returns your user ID.
func (w *SWhoAmI) WhoAmI() (string, error) {
	return w.Settings.UsersID, nil
}
