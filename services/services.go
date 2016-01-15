package services

import (
	"fmt"

	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
)

// CmdServices lists the names of all services for an environment.
func CmdServices(envID, pod string, is IServices) error {
	svcs, err := is.List(envID, pod)
	if err != nil {
		return err
	}
	fmt.Println("NAME")
	for _, s := range *svcs {
		fmt.Printf("%s\n", s.Label)
	}
	return nil
}

func (s *SServices) List(envID, pod string) (*[]models.Service, error) {
	headers := httpclient.GetHeaders(s.Settings.APIKey, s.Settings.SessionToken, s.Settings.Version, pod)
	resp, statusCode, err := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/services", s.Settings.PaasHost, s.Settings.PaasHostVersion, envID), headers)
	if err != nil {
		return nil, err
	}
	var services []models.Service
	err = httpclient.ConvertResp(resp, statusCode, &services)
	if err != nil {
		return nil, err
	}
	return &services, nil
}

func (s *SServices) Retrieve(envID, svcID, pod string) (*models.Service, error) {
	headers := httpclient.GetHeaders(s.Settings.APIKey, s.Settings.SessionToken, s.Settings.Version, pod)
	resp, statusCode, err := httpclient.Get(fmt.Sprintf("%s%s/environments/%s/services/%s", s.Settings.PaasHost, s.Settings.PaasHostVersion, envID, svcID), headers)
	if err != nil {
		return nil, err
	}
	var service models.Service
	err = httpclient.ConvertResp(resp, statusCode, &service)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (s *SServices) RetrieveByLabel(label, envID, pod string) (*models.Service, error) {
	services, err := s.List(envID, pod)
	if err != nil {
		return nil, err
	}
	var service *models.Service
	for _, s := range *services {
		if s.Label == label {
			service = &s
			break
		}
	}
	return service, nil
}
