// client.go
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

//	"io/ioutil"
//	"crypto/tls"
//	"crypto/x509"

	"github.com/gorilla/websocket"
)

var done chan interface{}
var interrupt chan os.Signal

func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		log.Printf("Received: %s\n", msg)
	}
}

func main() {
	done = make(chan interface{}) // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	socketUrl := "ws://localhost:8080" + "/messaging"
/*
	socketUrl := "wss://localhost:8080" + "/messaging"
//	securewebsocket := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	roots := x509.NewCertPool()
	var rootPEM string = ""
	crt, errCrt := ioutil.ReadFile("./cert.crt")
	if(errCrt != nil) {
		log.Fatal("Failed to read root certificate CRT")
	}
	key, errKey := ioutil.ReadFile("./cert.key")
	if(errKey != nil) {
		log.Fatal("Failed to read root certificate Key")
	}
	rootPEM = string(crt) + "\n" + string(key)
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		log.Fatal("Failed to parse root certificate")
	}
	securewebsocket := websocket.Dialer{TLSClientConfig: &tls.Config{RootCAs: roots}}
	conn, _, err := securewebsocket.Dial(socketUrl, nil)
*/
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()
	go receiveHandler(conn)

	// Our main loop for the client
	// We send our relevant packets here
	for {
		select {
			case <-time.After(time.Duration(5) * time.Millisecond * 1000):
				// Send an echo packet every second
				fmt.Println("Sending message to server")
				err := conn.WriteMessage(websocket.TextMessage, []byte("Hello from Client ..."))
				if err != nil {
					log.Println("Error during writing to websocket:", err)
					return
				}
			case <-interrupt:
				// We received a SIGINT (Ctrl + C). Terminate gracefully...
				log.Println("Received SIGINT interrupt signal. Closing all pending connections")
				// Close our websocket connection
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("Error during closing websocket:", err)
					return
				}
				select {
					case <-done:
						log.Println("Receiver Channel Closed! Exiting....")
					case <-time.After(time.Duration(1) * time.Second):
						log.Println("Timeout in closing receiving channel. Exiting....")
				}
				return
		}
	}

}


