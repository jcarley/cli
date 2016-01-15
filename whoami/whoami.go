package whoami

import "fmt"

func CmdWhoAmI(w IWhoAmI) error {
	usersID, err := w.WhoAmI()
	if err != nil {
		return err
	}
	fmt.Printf("user ID = %s\n", usersID)
	return nil
}

// WhoAmI prints out your user ID. This can be used for adding users to
// environments via `catalyze adduser`, removing users from an environment
// via `catalyze rmuser`, when contacting Catalyze Support, etc.
func (w *SWhoAmI) WhoAmI() (string, error) {
	return w.Settings.UsersID, nil
}
