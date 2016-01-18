package users

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "users",
	ShortHelp: "Manage users who have access to the associated environment",
	LongHelp:  "Manage users who have access to the associated environment",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			// TODO make this work with the core api orgs
			cmd.Command(AddSubCmd.Name, AddSubCmd.ShortHelp, AddSubCmd.CmdFunc(settings))
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(RmSubCmd.Name, RmSubCmd.ShortHelp, RmSubCmd.CmdFunc(settings))
		}
	},
}

var AddSubCmd = models.Command{
	Name:      "add",
	ShortHelp: "Grant access to the associated environment for the given user",
	LongHelp:  "Grant access to the associated environment for the given user",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			usersID := subCmd.StringArg("USER_ID", "", "The Users ID to give access to the associated environment")
			subCmd.Action = func() {
				err := CmdAdd(*usersID, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "USER_ID"
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all users who have access to the associated environment",
	LongHelp:  "List all users who have access to the associated environment",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				err := CmdList(settings.UsersID, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

var RmSubCmd = models.Command{
	Name:      "rm",
	ShortHelp: "Revoke access to the associated environment for the given user",
	LongHelp:  "Revoke access to the associated environment for the given user",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			usersID := subCmd.StringArg("USER_ID", "", "The Users ID to revoke access from for the associated environment")
			subCmd.Action = func() {
				err := CmdRm(*usersID, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "USER_ID"
		}
	},
}

// IUsers
type IUsers interface {
	Add(usersID string) error
	List() (*models.EnvironmentUsers, error)
	Rm(usersID string) error
}

// SUsers is a concrete implementation of IUsers
type SUsers struct {
	Settings *models.Settings
}

// New generates a new instance of IUsers
func New(settings *models.Settings) IUsers {
	return &SUsers{
		Settings: settings,
	}
}
