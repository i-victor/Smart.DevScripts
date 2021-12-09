
// GO Lang
// go build main.go
// on openbsd may need to: CGO_LDFLAGS_ALLOW='-Wl,-z,wxneeded|-Wl,-rpath-link,/usr/X11R6/lib' go build main.go

// noVNC Viewer
// (c) 2018-2021 unix-world.org
// License: BSD

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"net/url"
	"path/filepath"

	"flag"
	"os"
	"golang.org/x/net/websocket"

	"embed"

	"github.com/unix-world/smartgo/webview"
//	"github.com/zserge/webview"
)

//go:embed assets/*
var assets embed.FS

var uxmScriptVersion = "r.20211206.1537"

var targetAddr = flag.String("target", "", "vnc-host:vnc-port")

var bindTcpAddr = ""
var vncTcpAddr = ""

func handleWss(wsconn *websocket.Conn) {
	log.Println("WebSocket Connection from Client")
	conn, err := net.Dial("tcp", *targetAddr)
	if err != nil {
		log.Println(err)
		wsconn.Close()
	} else {
		wsconn.PayloadType = websocket.BinaryFrame
		go io.Copy(conn, wsconn)
		go io.Copy(wsconn, conn)
		select {}
	}
}

func bootHandshake(config *websocket.Config, r *http.Request) error {
	config.Protocol = []string{"binary"}
	r.Header.Set("Access-Control-Allow-Origin", "*")
	r.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
	return nil
}

func startServer() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	bindTcpAddr = ln.Addr().String()
	fmt.Println("Listening on TCP: " + bindTcpAddr)
	go func() {
		defer ln.Close()
		fmt.Println("Starting the Built-in HTTP Server ...")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if len(path) > 0 && path[0] == '/' {
				path = path[1:]
			}
			if path == "" {
				path = "vnc-auto.html"
			}
		//	if bs, err := Asset(path); err != nil { // old, using go-bindata
			if bs, err := assets.ReadFile("assets/" + path); err != nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.Header().Add("Z-No-Vnc-Host-Port-Config", url.QueryEscape(ln.Addr().String()))
				w.Header().Add("Content-Type", mime.TypeByExtension(filepath.Ext(path)))
				io.Copy(w, bytes.NewBuffer(bs))
			}
		})
		log.Fatal(http.Serve(ln, nil))
	}()
	time.Sleep(time.Duration(1) * time.Second)
	go func() {
		fmt.Println("Starting the Built-in WebSocket Server ...")
		mux := websocket.Server { Handshake: bootHandshake, Handler: handleWss }
		http.Handle("/websockify", mux)
	}()
	time.Sleep(time.Duration(1) * time.Second)
	return "http://" + ln.Addr().String()
}

func wkExecRPC(w webview.WebView, msg string) {
	w.Eval(fmt.Sprintf("smartGoWkRPCCallExec(%s)", string(msg)))
}

func getTcpAddr() string {
	return bindTcpAddr
}

func getVncAddr() string {
	return vncTcpAddr
}

func handleRPC(w webview.WebView, data string) {
	jdata := struct {
		Cmd string `json:"cmd"`
		Msg string `json:"msg"`
		Val string `json:"val"`
	}{}
	if err := json.Unmarshal([]byte(data), &jdata); err != nil {
		fmt.Println("JSON-RPC ERROR !")
		fmt.Println(err)
		return
	}
	switch jdata.Cmd {
		case "close":
			w.Terminate()
		case "fullscreen":
			w.SetFullscreen(true)
		case "unfullscreen":
			w.SetFullscreen(false)
		case "init":
			//wkExecRPC(w, "alert('Init')")
			// do nothing
		case "novncsettings":
			wkExecRPC(w, "parseGoSettings('" + getTcpAddr() + "', '" + getVncAddr() + "')")
	}
}

func main() {

	flag.Parse()
	if *targetAddr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	vncTcpAddr = *targetAddr

	fmt.Println("##### GO.noVNC [ " + uxmScriptVersion + " ] :: (c) 2018-2021 unix-world.org")

	url := startServer()

	fmt.Println("Starting WebView ...")
	time.Sleep(time.Duration(1) * time.Second)
	w := webview.New(webview.Settings{
		Width:  1366,
		Height: 795,
		Title:  "noVNC",
		URL:    url,
		ExternalInvokeCallback: handleRPC,
	})

	fmt.Println("Done.")

	defer w.Exit()
	w.Run()

	fmt.Println("Closing Application ... Done.")

}

// #END
