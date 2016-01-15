package logout

func CmdLogout(il ILogout) error {
	return il.Logout()
}

// Logout clears the stored user information from the local machine. This does
// not remove environment data.
func (l *SLogout) Logout() error {
	l.Settings.SessionToken = ""
	l.Settings.UsersID = ""
	return nil
}
