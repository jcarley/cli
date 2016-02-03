package keys

import (
	"encoding/json"
	"fmt"

	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

// List all keys belonging to the auth'd user
func List(settings *models.Settings) ([]models.UserKey, error) {
	headers := httpclient.GetHeaders(settings.SessionToken, settings.Version, settings.Pod)
	resp, status, err := httpclient.Get(nil, fmt.Sprintf("%s%s/keys", settings.AuthHost, settings.AuthHostVersion), headers)
	if err != nil {
		return nil, err
	}

	keys := []models.UserKey{}
	err = httpclient.ConvertResp(resp, status, &keys)
	return keys, err
}

// Add adds a new public key to the authenticated user's account
func Add(settings *models.Settings, name string, publicKey string) error {
	body, err := json.Marshal(struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	}{
		Key:  publicKey,
		Name: name,
	})
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(settings.SessionToken, settings.Version, settings.Pod)
	resp, status, err := httpclient.Post(body, fmt.Sprintf("%s%s/keys", settings.AuthHost, settings.AuthHostVersion), headers)
	if err != nil {
		return err
	}
	if httpclient.IsError(status) {
		return httpclient.ConvertError(resp, status)
	}
	return nil
}

// Remove removes a public key by name from the authenticated user's account, returning an error if unsuccessful.
func Remove(settings *models.Settings, name string) error {
	headers := httpclient.GetHeaders(settings.SessionToken, settings.Version, settings.Pod)
	resp, status, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/keys/%s", settings.AuthHost, settings.AuthHostVersion, name), headers)
	if err != nil {
		return err
	}
	if httpclient.IsError(status) {
		return httpclient.ConvertError(resp, status)
	}
	return nil
}
