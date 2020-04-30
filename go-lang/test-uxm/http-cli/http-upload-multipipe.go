
// multipart upload with io.Pipe
// curl --insecure -X POST -H "Content-Type: multipart/form-data" -F 'webdav_action=upf' -F 'file=@page.pdf' https://user:pass@127.0.0.1/sites/smart-framework/admin.php/page/cloud.files/~/uploads

package main

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"crypto/tls"
//	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cheggaaa/pb"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("Usage: pipeUp <filename>\n")
		os.Exit(1)
	}

	input, err := os.Open(os.Args[1])
	check(err)
	defer input.Close()

	stat, err := input.Stat()
	check(err)

	pipeOut, pipeIn := io.Pipe()
	fsize := stat.Size()
	bar := pb.New(int(fsize)).SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true

	writer := multipart.NewWriter(pipeIn)

	// auth and extra form fields
	var username string = "admin"
	var passwd string = "the-pass"
	extraFormFields := map[string]string{
		"webdav_action": "upf",
		"client": "golang",
	}

	// do the request concurrently
	var resp *http.Response
	done := make(chan error)
	go func() {

		// prepare request
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		req, err := http.NewRequest("POST", "https://127.0.0.1/sites/smart-framework/admin.php/page/cloud.files/~/uploads/", pipeOut)
		if username != "" {
			req.SetBasicAuth(username, passwd)
		}
		if err != nil {
			done <- err
			return
		}
		req.ContentLength = fsize // filesize
		req.ContentLength += 227  // multipart header excluding filename
		req.ContentLength += int64(len(filepath.Base(os.Args[1])))
		req.ContentLength += 5

		for k, v := range extraFormFields {
			req.ContentLength += int64(len(k))
			req.ContentLength += int64(len(v))
		}
		req.ContentLength += 218

		req.Header.Set("Content-Type", writer.FormDataContentType())
		log.Println("Created Request")
		bar.Start()

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			done <- err
			return
		}
		switch resp.StatusCode {
			case 200:
				bar.FinishPrint("Status: 200 OK")
				break
			case 201:
				bar.FinishPrint("Status: 201 ACCEPTED")
				break
			default:
			//	bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
			//	bodyString := string(bodyBytes)
				bar.FinishPrint("Upload Status: NOT OK: " + resp.Status) // + " :: " + bodyString)
		}

		done <- nil
	}()

	for k, v := range extraFormFields {
		writer.WriteField(k, v)
	}

	part, err := writer.CreateFormFile("file", filepath.Base(os.Args[1]))
	check(err)

	out := io.MultiWriter(part, bar)
	_, err = io.Copy(out, input)
	check(err)

	check(writer.Close())
	check(pipeIn.Close()) // need to close the pipe to

	check(<-done)

	bar.FinishPrint("Done.")

}

func check(err error) {
	_, file, line, _ := runtime.Caller(1)
	if err != nil {
		log.Fatalf("Fatal from <%s:%d>\nError:%s", file, line, err)
	}
}
