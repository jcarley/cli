package disassociate

import (
	"fmt"

	"github.com/catalyzeio/cli/config"
)

// Disassociate removes an existing association with the environment. The
// `catalyze` remote on the local github repo will *NOT* be removed.
func (d *SDisassociate) Disassociate() error {
	// DeleteBreadcrumb removes the environment from the settings.Environments
	// array for you
	config.DeleteBreadcrumb(d.Alias, d.Settings)
	fmt.Printf("WARNING: Your existing git remote *has not* been removed.\n\n")
	fmt.Println("Association cleared.")
	return nil
}
