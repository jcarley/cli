package console

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"golang.org/x/net/websocket"

	"github.com/catalyzeio/cli/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/services"
	"github.com/docker/docker/pkg/term"
)

func CmdConsole(svcName, command string, ic IConsole, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the name \"%s\"\n", svcName)
	}
	return ic.Open(command, service)
}

// Open opens a secure console to a code or database service. For code
// services, a command is required. This command is executed as root in the
// context of the application root directory. For database services, no command
// is needed - instead, the appropriate command for the database type is run.
// For example, for a postgres database, psql is run.
func (c *SConsole) Open(command string, service *models.Service) error {
	fmt.Printf("Opening console to %s (%s)\n", service.Name, service.ID)
	task, err := c.Request(command, service)
	if err != nil {
		return err
	}
	fmt.Print("Waiting for the console to be ready. This might take a minute.")

	jobID, err := c.Tasks.PollForConsole(task, service)
	if err != nil {
		return err
	}
	defer c.Destroy(jobID, service)
	creds, err := c.RetrieveTokens(jobID, service)
	if err != nil {
		return err
	}

	creds.URL = strings.Replace(creds.URL, "http", "ws", 1)
	fmt.Println("Connecting...")

	// BEGIN websocket impl
	config, _ := websocket.NewConfig(creds.URL, "ws://localhost:9443/")
	config.TlsConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	config.Header["X-Console-Token"] = []string{creds.Token}
	ws, err := websocket.DialConfig(config)
	if err != nil {
		return err
	}
	defer ws.Close()
	fmt.Println("Connection opened")

	stdin, stdout, _ := term.StdStreams()
	fdIn, isTermIn := term.GetFdInfo(stdin)
	if !isTermIn {
		return errors.New("StdIn is not a terminal")
	}
	oldState, err := term.SetRawTerminal(fdIn)
	if err != nil {
		return err
	}
	defer term.RestoreTerminal(fdIn, oldState)

	signal.Notify(make(chan os.Signal, 1), os.Interrupt)

	done := make(chan struct{}, 2)
	go readWS(ws, stdout, done)
	go readStdin(stdin, ws, done)

	<-done
	return nil
}

func (c *SConsole) Request(command string, service *models.Service) (*models.Task, error) {
	console := map[string]string{}
	if command != "" {
		console["command"] = command
	}
	b, err := json.Marshal(console)
	if err != nil {
		return nil, err
	}
	headers := httpclient.GetHeaders(c.Settings.SessionToken, c.Settings.Version, c.Settings.Pod)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/services/%s/console", c.Settings.PaasHost, c.Settings.PaasHostVersion, service.ID), headers)
	if err != nil {
		return nil, err
	}
	// TODO this is broken. The task is returned in the Location header as a route
	var m map[string]string
	err = httpclient.ConvertResp(resp, statusCode, &m)
	if err != nil {
		return nil, err
	}
	return &models.Task{
		ID: m["taskId"],
	}, nil
}

func (c *SConsole) RetrieveTokens(jobID string, service *models.Service) (*models.ConsoleCredentials, error) {
	tokenRequest := map[string]string{
		"serviceid": service.ID,
		"jobid":     jobID,
	}
	b, err := json.Marshal(tokenRequest)
	if err != nil {
		return nil, err
	}
	headers := httpclient.GetHeaders(c.Settings.SessionToken, c.Settings.Version, c.Settings.Pod)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/console/token", c.Settings.PaasHost, c.Settings.PaasHostVersion), headers)
	if err != nil {
		return nil, err
	}
	var credentials models.ConsoleCredentials
	err = httpclient.ConvertResp(resp, statusCode, &credentials)
	if err != nil {
		return nil, err
	}
	return &credentials, nil
}

func (c *SConsole) Destroy(jobID string, service *models.Service) error {
	headers := httpclient.GetHeaders(c.Settings.SessionToken, c.Settings.Version, c.Settings.Pod)
	resp, statusCode, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/jobs/%s", c.Settings.PaasHost, c.Settings.PaasHostVersion, jobID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}

// Reads incoming data from the websocket and forwards it to stdout.
func readWS(ws *websocket.Conn, t io.Writer, done chan struct{}) {
	_, err := io.Copy(t, ws)
	if err == io.EOF {
		fmt.Println("Connection closed")
	} else if err != nil {
		fmt.Printf("Error reading data from server: %s", err)
	}
	done <- struct{}{}
}

// Reads data from stdin and writes it to the websocket.
func readStdin(t io.ReadCloser, ws *websocket.Conn, done chan struct{}) {
	_, err := io.Copy(ws, t)
	if err == io.EOF {
		fmt.Println("Input closed")
	} else if err != nil {
		fmt.Printf("Error writing data to server: %s", err)
	}
	done <- struct{}{}
}
