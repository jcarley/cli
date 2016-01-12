package defaultcmd

import (
	"fmt"

	"github.com/catalyzeio/cli/config"
)

// Set sets the default environment. This environment must already be
// associated. Any commands run outside of a git directory will use the default
// environment for context.
func (d *SDefault) Set() error {
	var found bool
	for name := range d.Settings.Environments {
		if name == d.Alias {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("No environment with an alias of \"%s\" has been associated. Please run \"catalyze associate\" first\n", d.Alias)
	}
	d.Settings.Default = d.Alias
	config.SaveSettings(d.Settings)
	fmt.Printf("%s is now the default environment\n", d.Alias)
	return nil
}
