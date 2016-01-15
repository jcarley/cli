package vars

import (
	"fmt"
	"sort"

	"github.com/catalyzeio/cli/helpers"
)

func CmdList(iv IVars) error {
	envVars, err := iv.List()
	if err != nil {
		return err
	}
	var keys []string
	for k := range envVars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("%s=%s\n", key, envVars[key])
	}
	return nil
}

// List lists all environment variables.
func (v *SVars) List() (map[string]string, error) {
	return helpers.ListEnvVars(v.Settings), nil
}
