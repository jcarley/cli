package disassociate

import (
	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/config"
)

func CmdDisassociate(alias string, id IDisassociate) error {
	err := id.Disassociate(alias)
	if err != nil {
		return err
	}
	logrus.Warnln("Your existing git remote *has not* been removed. You must do this manually.")
	logrus.Println("Association cleared.")
	return nil
}

// Disassociate removes an existing association with the environment. The
// `datica` remote on the local github repo will *NOT* be removed.
func (d *SDisassociate) Disassociate(alias string) error {
	// DeleteBreadcrumb removes the environment from the settings.Environments
	// array for you
	return config.DeleteBreadcrumb(alias, d.Settings)
}
