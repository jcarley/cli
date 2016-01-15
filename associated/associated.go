package associated

import "fmt"

func CmdAssociated(ia IAssociated) error {
	return ia.Associated()
}

// Associated lists all currently associated environments.
func (a *SAssociated) Associated() error {
	for envAlias, env := range a.Settings.Environments {
		fmt.Printf(`%s:
    Environment ID:   %s
    Environment Name: %s
    Service ID:       %s
    Associated at:    %s
    Default:          %v
    Pod:              %s
`, envAlias, env.EnvironmentID, env.Name, env.ServiceID, env.Directory, a.Settings.Default == envAlias, env.Pod)
	}
	if len(a.Settings.Environments) == 0 {
		fmt.Println("No environments have been associated")
	}
	return nil
}
