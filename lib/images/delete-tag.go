package images

import (
	"errors"
	"fmt"
	"net/url"
)

// DeleteTag deletes a tag for an image.
func (d *SImages) DeleteTag(imageName, tagName string) error {
	headers := d.Settings.HTTPManager.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod, d.Settings.UsersID)
	resp, statusCode, err := d.Settings.HTTPManager.Delete(nil, fmt.Sprintf("%s%s/environments/%s/images/%s/tags/%s", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, url.PathEscape(url.PathEscape(imageName)), tagName), headers)
	if err != nil {
		return err
	}
	if statusCode >= 400 {
		converted, convertErr := d.Settings.HTTPManager.ConvertError(resp, statusCode)
		if convertErr != nil {
			return convertErr
		}
		if converted.Code == 98005 {
			return errors.New("Unable to delete tag - it is currently in use as a release by one or more services.")
		}
	}
	return d.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
