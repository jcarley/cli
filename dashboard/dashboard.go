package dashboard

import "github.com/skratchdot/open-golang/open"

// Open opens up the Catalyze Dashboard in the default browser
func (d *SDashboard) Open() error {
	return open.Run("https://dashboard.catalyze.io")
}
