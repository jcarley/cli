// +build windows

package update

import (
	"fmt"
)

func verifyExeDirWriteable() error {
	return fmt.Errorf("This function is not intended to be called on Windows.")
}
