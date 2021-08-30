
package main

// server.go
// r.20210819.0346 :: STABLE

import (
	"log"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	smart "github.com/unix-world/smartgo"
)

const (
	msgPeriod 		= 10 * time.Second
)

var targetAddr = flag.String("bind", "localhost:8887", "host:port (Ex: localhost:8887)")

var serverID string = "default"

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

var upgrader = websocket.Upgrader{
	ReadBufferSize:    16384,
	WriteBufferSize:   16384,
//	EnableCompression: true,
} // use default options

func broadcastMsg(conn *websocket.Conn, rAddr string) {
	for {
		msg, errMsg := ComposePakMessage("Hi from Server:" + serverID, smart.JsonEncode("Hi from the Server"))
		if(errMsg != "") {
			log.Println("[ERROR] Send Message to Client:", errMsg)
			conn.Close()
			return
		} else {
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("[ERROR] Send Message to Client / Writing to websocket Failed:", err)
				conn.Close()
				return
			} else {
				log.Println("[OK] Send Message to Client:", rAddr)
			}
		}
		time.Sleep(msgPeriod)
	}
}

func socketHandler(w http.ResponseWriter, r *http.Request) {

	// Upgrade our raw HTTP connection to a websocket based one
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // this is for ths js client connected from another origin ...
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] Connection Upgrade Failed:", err)
		return
	}
	conn.SetReadLimit(10 * 1000 * 1000) // 10 MB
	defer conn.Close()
	log.Println("New Connection to:", conn.LocalAddr(), "From:", r.RemoteAddr)

	// The event loop
	go broadcastMsg(conn, r.RemoteAddr)
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("[ERROR] Message Reading Failed:", err)
			break
		}
		log.Println("[INFO] Got New Message from Client:", conn.LocalAddr())
		if(messageType == websocket.TextMessage) {
			ok, errMsg := HandleMessage(serverID, message, conn)
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

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index Page")
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

	flag.Parse()
	var addr string = *targetAddr
	if(addr == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

	serverID = GenerateUUID()

	// Auth: Authenticate as for a plain http request. Upgrade to the websocket protocol after authentication.
	// Most web applications use a token stored in a cookie for authentication. That also works for the websocket endpoint.

	http.HandleFunc("/messaging", socketHandler)
	http.HandleFunc("/", home)
//	log.Fatal("[ERROR]", http.ListenAndServe("localhost:8887", nil))
	log.Fatal("[ERROR]", http.ListenAndServeTLS(addr, "./cert.crt", "./cert.key", nil))

}

// #END
