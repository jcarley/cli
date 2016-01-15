package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/updater"
)

func getClient() *http.Client {
	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		},
	}

	return &http.Client{
		Transport: tr,
	}
}

func GetHeaders(apiKey, sessionToken, version, pod string) map[string][]string {
	return map[string][]string{
		"Accept":        {"application/json"},
		"Content-Type":  {"application/json"},
		"X-Api-Key":     {apiKey},
		"Authorization": {fmt.Sprintf("Bearer %s", sessionToken)},
		"X-CLI-Version": {version},
		"X-Pod-ID":      {pod},
	}
}

// ConvertResp takes in a resp from one of the httpclient methods and
// checks if it is a successful request. If not, it is parsed as an error object
// and returned as an error. Otherwise it will be marshalled into the requested
// interface. ALWAYS PASS A POINTER INTO THIS METHOD. If you don't pass a struct
// pointer your original object will be nil or an empty struct.
func ConvertResp(b []byte, statusCode int, s interface{}) error {
	if statusCode < 200 || statusCode >= 300 {
		msg := ""
		var errs models.Error
		err := json.Unmarshal(b, &errs)
		if err == nil && errs.Title != "" && errs.Description != "" {
			msg = fmt.Sprintf("(%d) %s: %s\n", errs.Code, errs.Title, errs.Description)
		} else {
			msg = fmt.Sprintf("(%d) %s\n", statusCode, string(b))
		}
		return errors.New(msg)
	}
	return json.Unmarshal(b, s)
}

// Get performs a GET request
func Get(url string, headers map[string][]string) ([]byte, int, error) {
	return MakeRequest("GET", url, nil, headers)
}

// Post performs a POST request
func Post(body []byte, url string, headers map[string][]string) ([]byte, int, error) {
	reader := bytes.NewReader(body)
	return MakeRequest("POST", url, reader, headers)
}

// PutFile uploads a file
func PutFile(filepath string, url string, headers map[string][]string) ([]byte, int, error) {
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return nil, 0, err
	}
	info, _ := file.Stat()
	client := getClient()
	req, _ := http.NewRequest("PUT", url, file)
	req.ContentLength = info.Size()

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	return respBody, resp.StatusCode, nil
}

// Put performs a PUT request
func Put(body []byte, url string, headers map[string][]string) ([]byte, int, error) {
	reader := bytes.NewReader(body)
	return MakeRequest("PUT", url, reader, headers)
}

// Delete performs a DELETE request
func Delete(url string, headers map[string][]string) ([]byte, int, error) {
	return MakeRequest("DELETE", url, nil, headers)
}

// MakeRequest is a generic HTTP runner that performs a request and returns
// the result body as a byte array. It's up to the caller to transform them
// into an object.
func MakeRequest(method string, url string, body io.Reader, headers map[string][]string) ([]byte, int, error) {
	client := getClient()
	req, _ := http.NewRequest(method, url, body)
	req.Header = headers

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 412 {
		updater.AutoUpdater.ForcedUpgrade()
		return nil, 0, fmt.Errorf("A required update has been applied. Please re-run this command.")
	}
	return respBody, resp.StatusCode, nil
}
