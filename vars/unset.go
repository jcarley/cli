package vars

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
)

func CmdUnset(key string, iv IVars) error {
	err := iv.Unset(key)
	if err != nil {
		return err
	}
	fmt.Println("Unset.")
	return nil
}

// Unset deletes an environment variable. Any changes to environment variables
// will not take effect until the service is redeployed by pushing new code
// or via `catalyze redeploy`.
func (v *SVars) Unset(variable string) error {
	helpers.UnsetEnvVar(variable, v.Settings)
	return nil
}
