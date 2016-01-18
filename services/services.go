package services

import (
	"fmt"

	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
)

// CmdServices lists the names of all services for an environment.
func CmdServices(is IServices) error {
	svcs, err := is.List()
	if err != nil {
		return err
	}
	fmt.Println("NAME")
	for _, s := range *svcs {
		fmt.Printf("%s\n", s.Label)
	}
	return nil
}

func (s *SServices) List() (*[]models.Service, error) {
	headers := httpclient.GetHeaders(s.Settings.SessionToken, s.Settings.Version, s.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services", s.Settings.PaasHost, s.Settings.PaasHostVersion, s.Settings.EnvironmentID), headers)
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

func (s *SServices) Retrieve(svcID string) (*models.Service, error) {
	headers := httpclient.GetHeaders(s.Settings.SessionToken, s.Settings.Version, s.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s", s.Settings.PaasHost, s.Settings.PaasHostVersion, s.Settings.EnvironmentID, s.Settings.ServiceID), headers)
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

func (s *SServices) RetrieveByLabel(label string) (*models.Service, error) {
	services, err := s.List()
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
