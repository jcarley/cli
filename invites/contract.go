package invites

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/prompts"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "invites",
	ShortHelp: "Manage invitations for your environments",
	LongHelp:  "Manage invitations for your environments",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			// TODO make use of orgs in the core api to use this command
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(RmSubCmd.Name, RmSubCmd.ShortHelp, RmSubCmd.CmdFunc(settings))
			cmd.Command(SendSubCmd.Name, SendSubCmd.ShortHelp, SendSubCmd.CmdFunc(settings))
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all pending environment invitations",
	LongHelp:  "List all pending environment invitations",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				err := CmdList(settings.EnvironmentName, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

var RmSubCmd = models.Command{
	Name:      "rm",
	ShortHelp: "Remove a pending environment invitation",
	LongHelp:  "Remove a pending environment invitation",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			inviteID := subCmd.StringArg("INVITE_ID", "", "The ID of an invitation to remove")
			subCmd.Action = func() {
				err := CmdRm(*inviteID, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "INVITE_ID"
		}
	},
}

var SendSubCmd = models.Command{
	Name:      "send",
	ShortHelp: "Send an invite to a user by email for the associated environment",
	LongHelp:  "Send an invite to a user by email for the associated environment",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			email := subCmd.StringArg("EMAIL", "", "The email of a user to invite to the associated environment. This user does not need to have a Catalyze account prior to sending the invitation")
			subCmd.Action = func() {
				err := CmdSend(*email, settings.EnvironmentName, New(settings), prompts.New())
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "EMAIL"
		}
	},
}

// IInvites
type IInvites interface {
	List() (*[]models.Invite, error)
	Rm(inviteID string) error
	Send(email string) error
}

// SInvites is a concrete implementation of IInvites
type SInvites struct {
	Settings *models.Settings
}

// New generates a new instance of IInvites
func New(settings *models.Settings) IInvites {
	return &SInvites{
		Settings: settings,
	}
}
