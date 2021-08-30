
package main

// client.go
// r.20210819.0346

import (
	"log"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crypto/tls"
//	"io/ioutil"
//	"crypto/x509"

	"github.com/gorilla/websocket"

	smart "github.com/unix-world/smartgo"
)

const (
	msgPeriod 		= 30 * time.Second
	reconnectPeriod	= 60 * time.Second
)

var targetAddr = flag.String("peer", "127.0.0.1:8887", "host:port (Ex: localhost:8887)")

func LogToConsoleWithColors() {
	//--
	smart.ClearPrintTerminal()
	//--
//	smart.LogToStdErr("DEBUG")
	smart.LogToConsole("DEBUG", true) // colored or not
//	smart.LogToFile("WARNING", "logs/", "json", true, true) // json | plain ; also on console ; colored or not
	//--
//	log.Println("[DATA] Data Debugging")
//	log.Println("[DEBUG] Debugging")
//	log.Println("[NOTICE] Notice")
//	log.Println("[WARNING] Warning")
//	log.Println("[ERROR] Error")
//	log.Println("[OK] OK")
//	log.Println("A log message, with no type") // aka [INFO]
	//--
} //END FUNCTION

func fatalError(logMessages ...interface{}) {
	//--
	log.Fatal("[ERROR] ", logMessages) // standard logger
	os.Exit(1)
	//--
} //END FUNCTION


var clientID string = "default"
var done chan interface{}
var interrupt chan os.Signal

func receiveHandler(connection *websocket.Conn) {
	defer smart.PanicHandler()
	defer close(done)
	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("[ERROR] Message Receive Failed:", err)
			return
		}
		log.Printf("[NOTICE] Received Message, Size = %d bytes\n", len(string(message)))
		if(messageType == websocket.TextMessage) {
			ok, errMsg := HandleMessage(clientID, message, nil)
			message = nil
			var retMsg string = "[OK] Message ..."
			if(ok != true) {
				if(errMsg == "") {
					errMsg = "[ERROR]: Unknown Error !!!"
				}
				retMsg = errMsg
			}
			log.Println(retMsg)
		} else {
			log.Println("[ERROR]: TextMessage is expected")
		}
	}
}

var connectedPeers = map[string]*websocket.Conn{}

func connectToPeer(addr string) {

	defer smart.PanicHandler()

	log.Println("[NOTICE] Connecting to Peer:", addr)

	socketUrl := "wss://" + addr + "/messaging"
	securewebsocket := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
/*
	roots := x509.NewCertPool()
	var rootPEM string = ""
	crt, errCrt := ioutil.ReadFile("./cert.crt")
	if(errCrt != nil) {
		log.Fatal("[ERROR] Failed to read root certificate CRT")
	}
	key, errKey := ioutil.ReadFile("./cert.key")
	if(errKey != nil) {
		log.Fatal("[ERROR] to read root certificate Key")
	}
	rootPEM = string(crt) + "\n" + string(key)
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		log.Fatal("[ERROR] Failed to parse root certificate")
	}
	securewebsocket := websocket.Dialer{TLSClientConfig: &tls.Config{RootCAs: roots}}
*/
	conn, _, err := securewebsocket.Dial(socketUrl, nil)
//	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Println("[ERROR] Cannot connect to Websocket Server: ", err)
		conn.Close()
		return
	}
	connectedPeers[addr] = conn
	defer conn.Close()
	go receiveHandler(conn)

	// Our main loop for the client
	// We send our relevant packets here
	for {
		select {
			case <-time.After(time.Duration(1) * msgPeriod):
				log.Println("[NOTICE] Sending message to server")
				msg, errMsg := ComposePakMessage("helloworld:" + clientID, smart.JsonEncode("This is a message for HelloWorld, from client"))
				if(errMsg != "") {
					log.Println("[ERROR]:", errMsg)
					conn.Close()
					delete(connectedPeers, addr)
					return
				} else {
					err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
					if err != nil {
						log.Println("[ERROR] Writing to websocket Failed:", err)
						conn.Close()
						delete(connectedPeers, addr)
						return
					}
				}
				msg = ""
				errMsg = ""
			case <-interrupt: // received a SIGINT (Ctrl + C). Terminate gracefully...
				log.Println("[NOTICE] Received SIGINT interrupt signal. Closing all pending connections")
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")) // Close our websocket connection
				if err != nil {
					log.Println("[ERROR] Closing websocket Failed:", err)
				}
			//	conn.Close()
				delete(connectedPeers, addr)
				return
				select {
					case <-done:
						log.Println("[NOTICE] Receiver Channel Closed...")
					case <-time.After(time.Duration(1) * time.Second):
						log.Println("[WARNING] Timeout in closing receiving channel...")
				}
				return
		}
	}

}

func main() {

	LogToConsoleWithColors()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("[INFO] CTRL+C Exit ...")
		os.Exit(1)
	}()

	done = make(chan interface{}) // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully
	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	flag.Parse()
	var addr string = *targetAddr
	if(addr == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

	clientID = GenerateUUID()
//	var pool = []string{"127.0.0.1:8887", "127.0.0.1:8888"}
	var pool = []string{"127.0.0.1:8887"}
	var loops int = 0;
	for {
		log.Println("[INFO] ... WATCHDOG ...")
	//	log.Println("[DATA] Connected Peers:", connectedPeers)
		for _, p := range pool {
			if _, exist := connectedPeers[p]; exist {
				log.Println("[OK] Peer is Connected to:", p)
			} else {
				if(loops > 0) {
					log.Println("[WARNING] Peer Not Connected to:", p)
				}
				go connectToPeer(p)
			}
		}
		loops++
		time.Sleep(reconnectPeriod)
	}

}

// #END
