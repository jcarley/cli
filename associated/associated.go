package associated

import (
	"fmt"

	"github.com/catalyzeio/cli/models"
)

// Associated lists all currently associated environments.
func Associated(settings *models.Settings) {
	for envAlias, env := range settings.Environments {
		fmt.Printf(`%s:
    Environment ID:   %s
    Environment Name: %s
    Service ID:       %s
    Associated at:    %s
    Default:          %v
    Pod:              %s
`, envAlias, env.EnvironmentID, env.Name, env.ServiceID, env.Directory, settings.Default == envAlias, env.Pod)
	}
	if len(settings.Environments) == 0 {
		fmt.Println("No environments have been associated")
	}
}
