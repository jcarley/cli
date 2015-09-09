package commands

import "github.com/skratchdot/open-golang/open"

// Dashboard opens up the Catalyze Dashboard in the default browser
func Dashboard() {
	open.Run("https://dashboard.catalyze.io")
}
