package sites

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/files"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/commands/ssl"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func CmdCreate(hostname, chainPath, privateKeyPath, serviceName string, wildcard, selfSigned bool, is ISites, issl ssl.ISSL, ifiles files.IFiles, iservices services.IServices) error {
	chainInfo, err := os.Stat(chainPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("A cert does not exist at path '%s'", chainPath)
	}
	privateKeyInfo, err := os.Stat(privateKeyPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("A private key does not exist at path '%s'", privateKeyPath)
	}
	err = issl.Verify(chainPath, privateKeyPath, hostname, selfSigned)
	if err != nil {
		return err
	}

	upstreamService, err := iservices.RetrieveByLabel(serviceName)
	if err != nil {
		return err
	}

	serviceProxy, err := iservices.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}

	chainServiceFile, err := ifiles.Create(serviceProxy.ID, chainPath, fmt.Sprintf("/etc/ssl/certs/%s", chainInfo.Name()), "0400")
	if err != nil {
		return err
	}

	privateKeyServiceFile, err := ifiles.Create(serviceProxy.ID, privateKeyPath, fmt.Sprintf("/etc/ssl/private/%s", privateKeyInfo.Name()), "0400")
	if err != nil {
		return err
	}

	site := models.Site{
		Name:                hostname,
		SSLCertFileID:       chainServiceFile.ID,
		SSLPrivateKeyFileID: privateKeyServiceFile.ID,
		Wildcard:            wildcard,
		ServiceName:         upstreamService.ID,
	}
	err = is.Create(serviceProxy.ID, &site)
	if err != nil {
		return err
	}
	logrus.Printf("Site created (ID = %s)\n", site.ID)
	return nil
}

func (s *SSites) Create(svcID string, site *models.Site) error {
	b, err := json.Marshal(site)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(s.Settings.SessionToken, s.Settings.Version, s.Settings.Pod)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/sites", s.Settings.PaasHost, s.Settings.PaasHostVersion, s.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return err
	}
	var createdSite models.Site
	err = httpclient.ConvertResp(resp, statusCode, &createdSite)
	if err != nil {
		return err
	}
	site.ID = createdSite.ID
	return nil
}
