package vars

import (
	"fmt"
	"strings"

	"github.com/catalyzeio/cli/helpers"
)

func CmdSet(variables []string, iv IVars) error {
	existingVars, err := iv.List()
	if err != nil {
		return err
	}
	envVarsMap := make(map[string]string, len(variables))
	for _, envVar := range variables {
		pieces := strings.SplitN(envVar, "=", 2)
		if len(pieces) != 2 {
			return fmt.Errorf("Invalid variable format. Expected <key>=<value> but got %s\n", envVar)
		}
		envVarsMap[pieces[0]] = pieces[1]
	}

	for key := range envVarsMap {
		if _, ok := existingVars[key]; ok {
			err := iv.Unset(key)
			if err != nil {
				return err
			}
		}
	}

	err = iv.Set(envVarsMap)
	if err != nil {
		return err
	}
	fmt.Println("Set. For these environment variables to take effect, you will need to redeploy your service with \"catalyze redeploy\"")
	return nil
}

// Set adds a new environment variables or updates the value of an existing
// environment variables. Any changes to environment variables will not take
// effect until the service is redeployed by pushing new code or via
// `catalyze redeploy`.
func (v *SVars) Set(envVarsMap map[string]string) error {
	helpers.SetEnvVars(envVarsMap, v.Settings)
	return nil
}
