package logout

import (
	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/lib/auth"
)

func CmdLogout(il ILogout, ia auth.IAuth) error {
	ia.Signout()
	logrus.Println("You have been logged out")
	return il.Clear()
}

// Clear clears the stored user information from the local machine. This does
// not remove environment data.
func (l *SLogout) Clear() error {
	l.Settings.SessionToken = ""
	l.Settings.UsersID = ""
	return nil
}
