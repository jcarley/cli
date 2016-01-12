package vars

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

// UnsetVar deletes an environment variable. Any changes to environment variables
// will not take effect until the service is redeployed by pushing new code
// or via `catalyze redeploy`.
func UnsetVar(variable string, settings *models.Settings) {
	helpers.SignIn(settings)
	helpers.UnsetEnvVar(variable, settings)
	fmt.Println("Unset.")
}
