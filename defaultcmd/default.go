package defaultcmd

import (
	"fmt"
	"os"

	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/models"
)

// SetDefault sets the default environment. This environment must already be
// associated. Any commands run outside of a git directory will use the default
// environment for context.
func SetDefault(alias string, settings *models.Settings) {
	var found bool
	for name := range settings.Environments {
		if name == alias {
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("No environment with an alias of \"%s\" has been associated. Please run \"catalyze associate\" first\n", alias)
		os.Exit(1)
	}
	settings.Default = alias
	config.SaveSettings(settings)
	fmt.Printf("%s is now the default environment\n", alias)
}
