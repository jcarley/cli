package version

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
)

func CmdVersion() error {
	versionString := fmt.Sprintf("version %s %s", config.VERSION, config.ArchString())
	logrus.Println(versionString)
	return nil
}
