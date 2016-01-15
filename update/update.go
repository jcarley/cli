package update

import (
	"fmt"

	"github.com/catalyzeio/cli/updater"
)

func CmdUpdate(iu IUpdate) error {
	fmt.Println("Checking for available updates...")
	needsUpdate, err := iu.Check()
	if err != nil {
		return err
	}
	if !needsUpdate {
		fmt.Println("You are already running the latest version of the Catalyze CLI")
		return nil
	}
	err = iu.Update()
	if err != nil {
		return err
	}
	fmt.Printf("Your CLI has been updated to version %s\n", updater.AutoUpdater.Info.Version)
	return nil
}

func (u *SUpdate) Check() (bool, error) {
	updater.AutoUpdater.FetchInfo()
	if updater.AutoUpdater.CurrentVersion == updater.AutoUpdater.Info.Version {
		return false, nil
	}
	return true, nil
}

// Update updates the  CLI if a new update is available.
func (u *SUpdate) Update() error {
	updater.AutoUpdater.ForcedUpgrade()
	return nil
}
