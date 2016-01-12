package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/catalyzeio/cli/config"
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

func setHeaders(req *http.Request, settings *models.Settings) {
	req.Header = http.Header{
		"Accept":        {"application/json"},
		"Content-Type":  {"application/json"},
		"X-Api-Key":     {config.APIKey},
		"Authorization": {fmt.Sprintf("Bearer %s", settings.SessionToken)},
		"X-CLI-Version": {config.VERSION},
		"X-Pod-ID":      {settings.Pod},
	}
}

// Get performs a GET request
func Get(url string, verify bool, settings *models.Settings) []byte {
	return MakeRequest("GET", url, nil, verify, settings)
}

// Post performs a POST request
func Post(body []byte, url string, verify bool, settings *models.Settings) []byte {
	reader := bytes.NewReader(body)
	return MakeRequest("POST", url, reader, verify, settings)
}

// PutFile uploads a file
func PutFile(filepath string, url string, verify bool, settings *models.Settings) []byte {
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	info, _ := file.Stat()
	client := getClient()
	req, _ := http.NewRequest("PUT", url, file)
	req.ContentLength = info.Size()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	if verify && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		fmt.Printf("%d %s\n", resp.StatusCode, string(respBody))
		os.Exit(1)
	}
	return respBody
}

// Put performs a PUT request
func Put(body []byte, url string, verify bool, settings *models.Settings) []byte {
	reader := bytes.NewReader(body)
	return MakeRequest("PUT", url, reader, verify, settings)
}

// Delete performs a DELETE request
func Delete(url string, verify bool, settings *models.Settings) []byte {
	return MakeRequest("DELETE", url, nil, verify, settings)
}

// MakeRequest is a generic HTTP runner that performs a request and returns
// the result body as a byte array. It's up to the caller to transform them
// into an object.
func MakeRequest(method string, url string, body io.Reader, verify bool, settings *models.Settings) []byte {
	client := getClient()
	req, _ := http.NewRequest(method, url, body)
	setHeaders(req, settings)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 412 {
		updater.AutoUpdater.ForcedUpgrade()
		fmt.Println("A required update has been applied. Please re-run this command.")
		os.Exit(1)
	} else if verify && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		var errors models.Errors
		err := json.Unmarshal(respBody, &errors)
		if err == nil && errors.ReportedErrors != nil && len(*errors.ReportedErrors) > 0 {
			for _, e := range *errors.ReportedErrors {
				fmt.Println(e.Message)
			}
		} else if err == nil && errors.Title != "" && errors.Description != "" {
			fmt.Printf("(%d) %s: %s\n", errors.Code, errors.Title, errors.Description)
		} else {
			fmt.Printf("%d %s\n", resp.StatusCode, string(respBody))
		}
		os.Exit(1)
	}
	return respBody
}
