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

var quitTries = 0

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
	fmt.Println("Connection opened")
	wsCh := make(chan bool)
	readCh := make(chan bool)

	// when we are ready to stop processing messages, send true through the channel
	stdin, stdout, _ := term.StdStreams()
	fdIn, isTermIn := term.GetFdInfo(stdin)
	if !isTermIn {
		panic(errors.New("StdIn is not a terminal"))
	}
	oldState, err := term.SetRawTerminal(fdIn)
	if err != nil {
		panic(err)
	}

	go webSocketDaemon(ws, &stdout, wsCh, readCh)

	cleanupHandler := make(chan os.Signal, 1)
	signal.Notify(cleanupHandler, os.Interrupt)
	go func() {
		for range cleanupHandler {
			quitTries++
			if quitTries > 1 {
				fmt.Println("Force closing")
				term.RestoreTerminal(fdIn, oldState)
				os.Exit(0)
			}
			go func() {
				// This is all really ugly, trapping ctrl-c. with the new console
				// gateway we will be able to remove this and properly shutdown when
				// we receive a close frame from the remote.
				fmt.Println("\nCleaning up")
				term.RestoreTerminal(fdIn, oldState)
				wsCh <- true
				readCh <- true
				// we need the destroy here because deferred statements do not get
				// called when you use os.Exit(). Again this will get cleaned up with
				// the new console gateway
				helpers.DestroyConsole(jobID, service.ID, settings)
				os.Exit(0)
			}()
		}
	}()

	defer term.RestoreTerminal(fdIn, oldState)
	quit := make(chan bool)
	tCh := make(chan []byte)
	go termDaemon(&stdin, tCh, quit, cleanupHandler)
passthrough:
	for {
		select {
		case msg := <-tCh:
			ws.Write(msg)
		case <-readCh:
			wsCh <- true
			quit <- true
			break passthrough
		}
	}
}

// handles outputting messages from the remote socket
func webSocketDaemon(ws *websocket.Conn, t *io.Writer, ch chan bool, readDone chan bool) {
	readCh := make(chan []byte)
	go readDaemon(ws, readCh, readDone)
	for {
		select {
		case msg := <-readCh:
			(*t).Write(msg)
		case <-ch:
			readDone <- true
			ws.Close()
			fmt.Println("Connection closed")
			return
		}
	}
}

// handles reading messages from the web socket
func readDaemon(ws *websocket.Conn, ch chan []byte, readDone chan bool) {
	for {
		msg := make([]byte, 1024)
		n, err := ws.Read(msg)
		if err != nil {
			if ws.IsServerConn() {
				fmt.Println("Error: " + err.Error())
			}
			readDone <- true
			return
		}
		ch <- msg[:n]
	}
}

// handles reading input from the terminal and passing it through the websocket
// connection. once ctrl+c or another quit command is issued, the chan is
// notified
func termDaemon(t *io.ReadCloser, ch chan []byte, quit chan bool, signalCh chan os.Signal) {
	reader := bufio.NewReader(*t)
	for {
		msg := make([]byte, 1024)
		select {
		case <-quit:
			return
		default:
			n, err := reader.Read(msg)
			// try and capture the ctrl-c byte (3)
			if n == 1 && msg[0] == 3 {
				signalCh <- os.Interrupt
				return
			}
			if err != nil {
				break
			}
			ch <- msg[:n]
		}
	}
}
