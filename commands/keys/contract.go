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
			cmd.Command(listSubCmd.Name, listSubCmd.ShortHelp, listSubCmd.CmdFunc(settings))
			cmd.Command(addSubCmd.Name, addSubCmd.ShortHelp, addSubCmd.CmdFunc(settings))
			cmd.Command(removeSubCmd.Name, removeSubCmd.ShortHelp, removeSubCmd.CmdFunc(settings))
			cmd.Command(setSubCmd.Name, setSubCmd.ShortHelp, setSubCmd.CmdFunc(settings))
		}
	},
}
