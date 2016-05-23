package releases

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd for keys
var Cmd = models.Command{
	Name:      "releases",
	ShortHelp: "Manage releases for code services",
	LongHelp:  "Manage releases for code services",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(RmSubCmd.Name, RmSubCmd.ShortHelp, RmSubCmd.CmdFunc(settings))
			cmd.Command(UpdateSubCmd.Name, UpdateSubCmd.ShortHelp, UpdateSubCmd.CmdFunc(settings))
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all releases for a given code service",
	LongHelp:  "List all releases for a given code service",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to list releases for")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdList(*serviceName, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err)
				}
			}
		}
	},
}

var RmSubCmd = models.Command{
	Name:      "rm",
	ShortHelp: "Remove a release from a code service",
	LongHelp:  "Remove a release from a code service",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to remove a release from")
			releaseName := cmd.StringArg("RELEASE_NAME", "", "The name of the release to remove")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRm(*serviceName, *releaseName, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err)
				}
			}
		}
	},
}

var UpdateSubCmd = models.Command{
	Name:      "update",
	ShortHelp: "Update a release from a code service",
	LongHelp:  "Update a release from a code service",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to update a release for")
			releaseName := cmd.StringArg("RELEASE_NAME", "", "The name of the release to update")
			notes := cmd.StringOpt("n notes", "", "The new notes to save on the release. If omitted, notes will be unchanged.")
			newReleaseName := cmd.StringOpt("r release", "", "The new name of the release. If omitted, the release name will be unchanged.")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdUpdate(*serviceName, *releaseName, *notes, *newReleaseName, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err)
				}
			}
			cmd.Spec = "SERVICE_NAME RELEASE_NAME [--notes] [--release]"
		}
	},
}

type IReleases interface {
	List(svcID string) (*[]models.Release, error)
	Rm(releaseName, svcID string) error
	Update(releaseName, svcID, notes, newReleaseName string) error
}

type SReleases struct {
	Settings *models.Settings
}

func New(settings *models.Settings) IReleases {
	return &SReleases{Settings: settings}
}
