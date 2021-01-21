
// GoLang Sample
// HTTP Client Upload Simple
// (c) 2020-2021 unix-world.org
// r.20210118.2155

// curl -X POST -H "Content-Type: application/octet-stream" --data-binary '@filename' http://127.0.0.1:5050/upload

package main

import (
	"log"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"io/ioutil"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Create("./upload/a.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	n, err := io.Copy(file, r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", n)))
}

func main() {
	http.HandleFunc("/upload", uploadHandler)

	go func() {
		time.Sleep(time.Second * 1)
		upload()
	}()

	http.ListenAndServe(":5050", nil)
}

func upload() {
	file, err := os.Open("./a.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	req, err := http.Post("http://127.0.0.1:5050/upload", "binary/octet-stream", file)
	req.Close = true
	if err != nil {
		log.Fatal(err)
		return
	}
	defer req.Body.Close()
	message, _ := ioutil.ReadAll(req.Body)
	fmt.Printf(string(message))
}
