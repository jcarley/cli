package httpclient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/catalyzeio/catalyze/models"
)

// the Dashboard API Key
const APIKey = "32a384f5-5d11-4214-812e-b35ced9af4d7"

func getClient() *http.Client {
	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
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
		"X-Api-Key":     {APIKey},
		"Authorization": {fmt.Sprintf("Bearer %s", settings.SessionToken)},
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
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return MakeRequest("PUT", url, file, verify, settings)
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
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	if verify && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		fmt.Println(fmt.Errorf("%d %s", resp.StatusCode, string(respBody)).Error())
		os.Exit(1)
	}
	return respBody
}
