package version

import (
	"fmt"
	"runtime"

	"github.com/catalyzeio/cli/config"
)

func Version() {
	archString := "other"
	switch runtime.GOARCH {
	case "386":
		archString = "32-bit"
	case "amd64":
		archString = "64-bit"
	case "arm":
		archString = "arm"
	}
	fmt.Printf("version %s %s\n", config.VERSION, archString)
}
