package update

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/bugsnag/osext"
	"github.com/catalyzeio/cli/lib/updater"
)

func CmdUpdate(iu IUpdate) error {
	needsUpdate, err := iu.Check()
	if err != nil {
		return err
	}
	// check if we can overwrite exe
	if needsUpdate && (runtime.GOOS == "linux" || runtime.GOOS == "darwin") {
		exe, err := osext.Executable()
		if err != nil {
			return exeGenericError()
		}
		exeDir := filepath.Dir(exe)
		f, err := os.Open(exeDir)
		if err != nil {
			return exeGenericError()
		}
		info, err := f.Stat()
		if err != nil {
			return exeGenericError()
		}
		usr, err := user.Current()
		if err != nil {
			return exeGenericError()
		}
		ownerId := strconv.FormatUint(uint64(info.Sys().(*syscall.Stat_t).Uid), 10)
		groupId := strconv.FormatUint(uint64(info.Sys().(*syscall.Stat_t).Gid), 10)
		// permissions are the 9 least significant bits
		mode := uint32(info.Mode()) & uint32((1<<9)-1)
		// Ex: 7   7   7
		//    111 111 111
		//         ^   ^
		groupWriteAble := uint32(1 << 4)
		globalWriteAble := uint32(1 << 1)
		// if user doesn't own the directory, and the directory isn't globally writeable,
		// and it is false that the directory is group writeable and the user is part of
		// the group, then we can't update.
		if ownerId != usr.Uid &&
			(mode&globalWriteAble) != globalWriteAble &&
			!((mode&groupWriteAble) == groupWriteAble && groupId == usr.Gid) {
			return fmt.Errorf("Your CLI cannot update, because your user cannot directly write to the directory it is in, \"%s\". Run the update command as the owner of the directory, or do the update manually.", exeDir)

		}
	}
	logrus.Println("Checking for available updates...")
	if !needsUpdate {
		logrus.Println("You are already running the latest version of the Catalyze CLI")
		return nil
	}
	err = iu.Update()
	if err != nil {
		return err
	}
	logrus.Printf("Your CLI has been updated to version %s", updater.AutoUpdater.Info.Version)
	return nil
}

func exeGenericError() error {
	return fmt.Errorf("There was an error trying to find where your CLI is on your system. You may need to manually update your CLI")
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
