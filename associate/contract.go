package associate

import (
	"fmt"
	"os"

	"github.com/catalyzeio/cli/environments"
	"github.com/catalyzeio/cli/git"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/pods"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "associate",
	ShortHelp: "Associates an environment",
	LongHelp:  "Associates an environment and service using the given alias to your local machine. For all further commands, the alias will be used instead. If no alias is given, the environment name is used.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			envName := cmd.StringArg("ENV_NAME", "", "The name of your environment")
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the primary code service to associate with this environment (i.e. 'app01')")
			alias := cmd.StringOpt("a alias", "", "A shorter name to reference your environment by for local commands")
			remote := cmd.StringOpt("r remote", "catalyze", "The name of the remote")
			defaultEnv := cmd.BoolOpt("d default", false, "Specifies whether or not the associated environment will be the default")
			cmd.Action = func() {
				//settings := r.GetSettings(false, false, *givenEnvName, *givenSvcName, baasHost, paasHost, *username, *password)
				// TODO this should be checked globablly and not just here
				if settings.Pods == nil || len(*settings.Pods) == 0 {
					p, err := pods.New(settings).List()
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}
					settings.Pods = p
					fmt.Println(settings.Pods)
				}
				err := CmdAssociate(*envName, *serviceName, *alias, *remote, *defaultEnv, New(settings), git.New(), environments.New(settings))
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			cmd.Spec = "ENV_NAME SERVICE_NAME [-a] [-r] [-d]"
		}
	},
}

// interfaces are the API calls
type IAssociate interface {
	Associate(name, remote string, defaultEnv bool, env *models.Environment, chosenService *models.Service) error
}

// SAssociate is a concrete implementation of IAssociate
type SAssociate struct {
	Settings *models.Settings
}

func New(settings *models.Settings) IAssociate {
	return &SAssociate{
		Settings: settings,
	}
}
