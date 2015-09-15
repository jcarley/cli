package commands

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"golang.org/x/net/websocket"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
	"github.com/docker/docker/pkg/term"
)

// Console opens a secure console to a code or database service. For code
// services, a command is required. This command is executed as root in the
// context of the application root directory. For database services, no command
// is needed - instead, the appropriate command for the database type is run.
// For example, for a postgres database, psql is run.
func Console(serviceLabel string, command string, settings *models.Settings) {
	helpers.SignIn(settings)
	service := helpers.RetrieveServiceByLabel(serviceLabel, settings)
	if service == nil {
		fmt.Printf("Could not find a service with the label \"%s\"\n", serviceLabel)
		os.Exit(1)
	}
	fmt.Printf("Opening console to %s (%s)\n", serviceLabel, service.ID)
	task := helpers.RequestConsole(command, service.ID, settings)
	fmt.Print("Waiting for the console to be ready. This might take a minute.")

	ch := make(chan string, 1)
	go helpers.PollConsoleJob(task.ID, service.ID, ch, settings)
	jobID := <-ch
	defer helpers.DestroyConsole(jobID, service.ID, settings)
	creds := helpers.RetrieveConsoleTokens(jobID, service.ID, settings)

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
		panic(err)
	}
	defer ws.Close()
	fmt.Println("Connection opened")

	stdin, stdout, _ := term.StdStreams()
	fdIn, isTermIn := term.GetFdInfo(stdin)
	if !isTermIn {
		panic(errors.New("StdIn is not a terminal"))
	}
	oldState, err := term.SetRawTerminal(fdIn)
	if err != nil {
		panic(err)
	}

	done := make(chan bool)
	msgCh := make(chan []byte, 2)
	go webSocketDaemon(ws, &stdout, done, msgCh)

	signal.Notify(make(chan os.Signal, 1), os.Interrupt)

	defer term.RestoreTerminal(fdIn, oldState)
	go termDaemon(&stdin, ws)
	<-done
}

// handles setting up a read daemon and outputting messages from the remote
// socket. If a websocket.CloseFrame is read, then the websocket connection
// is properly closed.
func webSocketDaemon(ws *websocket.Conn, t *io.Writer, done chan bool, msgCh chan []byte) {
	go readDaemon(ws, msgCh)
	for {
		msg := <-msgCh
		if len(msg) == 1 && msg[0] == websocket.CloseFrame {
			ws.Close()
			fmt.Println("Connection closed")
			done <- true
			return
		}
		(*t).Write(msg)
	}
}

// handles reading messages from the web socket and passing them through the
// given chan. If an error occurs, the websocket.CloseFrame signal is sent
// through the chan.
func readDaemon(ws *websocket.Conn, msgCh chan []byte) {
	for {
		msg := make([]byte, 1024)
		n, err := ws.Read(msg)
		if err != nil {
			msgCh <- []byte{websocket.CloseFrame}
			return
		}
		msgCh <- msg[:n]
	}
}

// handles reading input from the terminal and passing it through the websocket
// connection.
func termDaemon(t *io.ReadCloser, ws *websocket.Conn) {
	reader := bufio.NewReader(*t)
	for {
		msg := make([]byte, 1024)
		n, err := reader.Read(msg)
		if err != nil {
			return
		}
		ws.Write(msg[:n])
	}
}
