package disassociate

import (
	"fmt"

	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/models"
)

// Disassociate removes an existing association with the environment. The
// `catalyze` remote on the local github repo will *NOT* be removed.
func Disassociate(alias string, settings *models.Settings) {
	// DeleteBreadcrumb removes the environment from the settings.Environments
	// array for you
	config.DeleteBreadcrumb(alias, settings)
	fmt.Printf("WARNING: Your existing git remote *has not* been removed.\n\n")
	fmt.Println("Association cleared.")
}
