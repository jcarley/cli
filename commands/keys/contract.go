package keys

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/deploykeys"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd for keys
var Cmd = models.Command{
	Name:      "keys",
	ShortHelp: "Tasks for SSH keys",
	LongHelp: "The `keys` command gives access to SSH key management for your user account. " +
		"SSH keys can be used for authentication and pushing code to the Catalyze platform. " +
		"Any SSH keys added to your user account should not be shared but be treated as private SSH keys. " +
		"Any SSH key uploaded to your user account will be able to be used with all code services and environments that you have access to. " +
		"The keys command can not be run directly but has sub commands.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(AddSubCmd.Name, AddSubCmd.ShortHelp, AddSubCmd.CmdFunc(settings))
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(RemoveSubCmd.Name, RemoveSubCmd.ShortHelp, RemoveSubCmd.CmdFunc(settings))
			cmd.Command(SetSubCmd.Name, SetSubCmd.ShortHelp, SetSubCmd.CmdFunc(settings))
		}
	},
}

var AddSubCmd = models.Command{
	Name:      "add",
	ShortHelp: "Add a public key",
	LongHelp: "`keys add` allows you to add a new SSH key to your user account. " +
		"SSH keys added to your user account should be private and not shared with others. " +
		"SSH keys can be used for authentication (as opposed to the traditional username and password) as well as pushing code to an environment's code services. " +
		"Please note, you must specify the path to the public key file and not the private key. " +
		"All SSH keys should be in either OpenSSH RSA format or PEM format. Here is a sample command\n\n" +
		"```catalyze keys add my_prod_key ~/.ssh/prod_rsa.pub```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			name := cmd.StringArg("NAME", "", "The name for the new key, for your own purposes")
			path := cmd.StringArg("PUBLIC_KEY_PATH", "", "Relative path to the public key file")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdAdd(*name, *path, New(settings), deploykeys.New(settings))
				if err != nil {
					logrus.Fatal(err)
				}
			}
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List your public keys",
	LongHelp: "`keys list` lists all public keys by name that have been uploaded to your user account including the key's fingerprint in SHA256 format. " +
		"Here is a sample command\n\n" +
		"```catalyze keys list```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdList(New(settings), deploykeys.New(settings))
				if err != nil {
					logrus.Fatal(err)
				}
			}
		}
	},
}

var RemoveSubCmd = models.Command{
	Name:      "rm",
	ShortHelp: "Remove a public key",
	LongHelp: "`keys rm` allows you to remove an SSH key previously uploaded to your account. " +
		"The name of the key can be found by using the [keys list](#keys-list) command. Here is a sample command\n\n" +
		"```catalyze keys rm my_prod_key```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			name := cmd.StringArg("NAME", "", "The name of the key to remove.")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRemove(*name, settings.PrivateKeyPath, New(settings), deploykeys.New(settings))
				if err != nil {
					logrus.Fatal(err)
				}
			}
		}
	},
}

var SetSubCmd = models.Command{
	Name:      "set",
	ShortHelp: "Set your auth key",
	LongHelp: "`keys set` allows the CLI to use an SSH key for authentication instead of the traditional username and password combination. " +
		"This can be useful for automation or where shared workstations are involved. " +
		"Please note that you must pass in the path to the private key and not the public key. " +
		"The given key must already be added to your account by using the [keys add](#keys-add) command. " +
		"Here is a sample command\n\n" +
		"```catalyze keys set ~/.ssh/my_key```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			path := cmd.StringArg("PRIVATE_KEY_PATH", "", "Relative path to the private key file")
			cmd.Action = func() {
				err := CmdSet(*path, settings)
				if err != nil {
					logrus.Fatal(err)
				}
			}
		}
	},
}

type IKeys interface {
	List() (*[]models.UserKey, error)
	Add(name, publicKey string) error
	Remove(name string) error
}

type SKeys struct {
	Settings *models.Settings
}

func New(settings *models.Settings) IKeys {
	return &SKeys{Settings: settings}
}
