
package main

// go-bindata -o assets.go -prefix assets assets/...
// go build main.go assets.go

import (
	"bytes"
//	"strings"
	"encoding/json"
	"fmt"
//	"time"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"path/filepath"

	"github.com/zserge/webview"
)

func startServer() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if len(path) > 0 && path[0] == '/' {
				path = path[1:]
			}
			if path == "" {
				path = "index.html"
			}
			if bs, err := Asset(path); err != nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.Header().Add("Content-Type", mime.TypeByExtension(filepath.Ext(path)))
				io.Copy(w, bytes.NewBuffer(bs))
			}
		})
		log.Fatal(http.Serve(ln, nil))
	}()
	return "http://" + ln.Addr().String()
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
			// do nothing
		case "test":
			fmt.Println(jdata.Msg + ": #" + jdata.Val)
	}
}

func main() {
	fmt.Println("Starting ...")
	url := startServer()
	fmt.Println("Done.")
	w := webview.New(webview.Settings{
		Width:  960,
		Height: 720,
		Title:  "Charts",
		URL:    url,
		ExternalInvokeCallback: handleRPC,
	})
	defer w.Exit()
	w.Run()
}

// #END
