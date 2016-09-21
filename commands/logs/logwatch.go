package logs

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

const (
	writeTimeout = 5 * time.Second
)

type LogMessage struct {
	Message   string `json:"message"`
	Timestamp string `json:"@timestamp"`
	Source    string `json:"source"`
}

func (l *SLogs) Watch(queryString, domain, sessionToken string) error {
	if queryString == "*" {
		queryString = ""
	}
	query, err := regexp.Compile(queryString)
	if err != nil {
		return err
	}
	logrus.Println("Streaming logs...")
	dialer := &websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	headers := http.Header{"Cookie": {"sessionToken=" + url.QueryEscape(sessionToken)}}
	urlString := fmt.Sprintf("wss://%s/stream", domain)
	c, _, err := dialer.Dial(urlString, headers)
	if err != nil {
		return err
	}
	defer c.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{}, 2)
	go func() {
		<-interrupt
		done <- struct{}{}
	}()
	go readWS(c, query, done)
	<-done
	logrus.Println("Disconnected")
	return nil
}

// Reads incoming data from the websocket and forwards it to stdout.
func readWS(ws *websocket.Conn, query *regexp.Regexp, done chan struct{}) {
	ws.SetPingHandler(func(string) error {
		ws.SetWriteDeadline(time.Now().Add(writeTimeout))
		return ws.WriteMessage(websocket.PongMessage, []byte{})
	})
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			done <- struct{}{}
			return
		}
		var log LogMessage
		err = json.Unmarshal(msg, &log)
		if err == nil {
			if query == nil || query.MatchString(log.Message) {
				logrus.Printf("%s - %s", log.Timestamp, log.Message)
			}
		} else {
			logrus.StandardLogger().Out.Write(msg)
		}
	}
}
