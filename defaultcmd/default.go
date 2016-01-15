package defaultcmd

import (
	"fmt"

	"github.com/catalyzeio/cli/config"
)

func CmdDefault(alias string, id IDefault) error {
	err := id.Set(alias)
	if err != nil {
		return err
	}
	fmt.Printf("%s is now the default environment\n", alias)
	return nil
}

// Set sets the default environment. This environment must already be
// associated. Any commands run outside of a git directory will use the default
// environment for context.
func (d *SDefault) Set(alias string) error {
	var found bool
	for name := range d.Settings.Environments {
		if name == alias {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("No environment with an alias of \"%s\" has been associated. Please run \"catalyze associate\" first\n", alias)
	}
	d.Settings.Default = alias
	config.SaveSettings(d.Settings)
	return nil
}
