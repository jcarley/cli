package domain

import (
	"errors"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/commands/sites"
)

// CmdDomain prints out the namespace plus domain of the given environment
func CmdDomain(envID string, ie environments.IEnvironments, is services.IServices, isites sites.ISites) error {
	env, err := ie.Retrieve(envID)
	if err != nil {
		return err
	}
	serviceProxy, err := is.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}
	sites, err := isites.List(serviceProxy.ID)
	if err != nil {
		return err
	}
	domain := ""
	for _, site := range *sites {
		if strings.HasPrefix(site.Name, env.Namespace) {
			domain = site.Name
		}
	}
	if domain == "" {
		return errors.New("Could not determine the temporary domain name of your environment")
	}
	logrus.Println(domain)
	return nil
}
