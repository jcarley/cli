package clear

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/models"
)

func CmdClear(privateKey, session, environments, defaultEnv, pods bool, settings *models.Settings) error {
	if defaultEnv {
		logrus.Warnln("The \"--default\" flag has been deprecated! It will be removed in a future version.")
	}
	if privateKey {
		settings.PrivateKeyPath = ""
	}
	if session {
		settings.SessionToken = ""
		settings.UsersID = ""
	}
	if environments {
		settings.Environments = map[string]models.AssociatedEnv{}
	}
	if defaultEnv {
		settings.Default = ""
	}
	if pods {
		settings.Pods = &[]models.Pod{}
	}
	config.SaveSettings(settings)
	if !privateKey && !session && !environments && !defaultEnv && !pods {
		logrus.Println("No settings were specified. To see available options, run \"catalyze clear --help\"")
	} else {
		logrus.Println("All specified settings have been cleared")
	}
	return nil
}
