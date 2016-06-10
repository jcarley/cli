// +build !windows

package update

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/bugsnag/osext"
)

func verifyExeDirWriteable() error {
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
	ownerID := strconv.FormatUint(uint64(info.Sys().(*syscall.Stat_t).Uid), 10)
	groupID := strconv.FormatUint(uint64(info.Sys().(*syscall.Stat_t).Gid), 10)
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
	if ownerID != usr.Uid &&
		(mode&globalWriteAble) != globalWriteAble &&
		!((mode&groupWriteAble) == groupWriteAble && groupID == usr.Gid) {
		return fmt.Errorf("Your CLI cannot update, because your user cannot directly write to the directory it is in, \"%s\". Run the update command as the owner of the directory, or do the update manually.", exeDir)
	}
	return nil
}
