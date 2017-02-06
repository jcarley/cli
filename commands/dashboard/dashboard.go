package dashboard

import (
	"github.com/daticahealth/cli/config"
	"github.com/skratchdot/open-golang/open"
)

func CmdDashboard(id IDashboard) error {
	return id.Open()
}

// Open opens up the Datica Dashboard in the default browser
func (d *SDashboard) Open() error {
	return open.Run(config.AccountsHost)
}
