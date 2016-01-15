package rake

import (
	"fmt"

	"github.com/catalyzeio/cli/helpers"
)

func CmdRake(taskName string, ir IRake) error {
	fmt.Printf("Executing Rake task: %s\n", taskName)
	err := ir.Run(taskName)
	if err != nil {
		return err
	}
	fmt.Println("Rake task output viewable in your logging server")
	return nil
}

// Rake executes a rake task. This is only applicable for ruby-based
// applications.
func (r *SRake) Run(taskName string) error {
	helpers.InitiateRakeTask(taskName, r.Settings)
	return nil
}
