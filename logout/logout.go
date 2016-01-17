package logout

import "github.com/catalyzeio/cli/auth"

func CmdLogout(il ILogout, ia auth.IAuth) error {
	err := ia.Signout()
	if err != nil {
		return err
	}
	return il.Clear()
}

// Clear clears the stored user information from the local machine. This does
// not remove environment data.
func (l *SLogout) Clear() error {
	l.Settings.SessionToken = ""
	l.Settings.UsersID = ""
	return nil
}
