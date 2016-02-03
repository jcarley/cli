package files

import (
	"fmt"

	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
)

// CmdRm removes a service file by its name.
func CmdRm(svcName, fileName string, ifiles IFiles, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the name \"%s\"", svcName)
	}
	file, err := ifiles.Retrieve(fileName, service.ID)
	if err != nil {
		return err
	}
	return ifiles.Rm(file.ID, service.ID)
}

func (f *SFiles) Rm(fileID int, svcID string) error {
	headers := httpclient.GetHeaders(f.Settings.SessionToken, f.Settings.Version, f.Settings.Pod)
	resp, statusCode, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/files/%d", f.Settings.PaasHost, f.Settings.PaasHostVersion, f.Settings.EnvironmentID, svcID, fileID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
