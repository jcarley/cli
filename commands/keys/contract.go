package keys

import (
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd for keys
var Cmd = models.Command{
	Name:      "keys",
	ShortHelp: "Tasks for SSH keys",
	LongHelp:  "Tasks for your own SSH keys",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(AddSubCmd.Name, AddSubCmd.ShortHelp, AddSubCmd.CmdFunc(settings))
			cmd.Command(RemoveSubCmd.Name, RemoveSubCmd.ShortHelp, RemoveSubCmd.CmdFunc(settings))
			cmd.Command(SetSubCmd.Name, SetSubCmd.ShortHelp, SetSubCmd.CmdFunc(settings))
		}
	},
}
