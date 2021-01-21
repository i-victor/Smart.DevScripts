
// GoLang Sample
// HTTP Client POST : file (multi-part form) | json (simple form) | text (simple form)
// (c) 2020-2021 unix-world.org
// r.20210118.2155

package main

import (
	"os"
	"log"
	"fmt"
	"bytes"
	"strings"
	"strconv"
	"path/filepath"
	"io"
	"mime/multipart"
	"net/http"
	"crypto/tls"
	"time"

	"github.com/cheggaaa/pb"
)


const ( // all tests can be performed with: webdav-server.go which also serves plain HTTP(S) (on /)
	THE_URL = "http://127.0.0.1:13080/" 		// Ex(text|json|file): `http://127.0.0.1:13080/` or `https://127.0.0.1:13443/`
	THE_AUTH_USERNAME = ""						// leave empty if no auth required or fill the auth username
	THE_AUTH_PASSWORD = ""						// leave empty if no auth required or fill the auth password
)


func postRequest(username string, passwd string, url string, fName string, data io.Reader, datalen int64) int {
	//--
	if(fName == "") {
		log.Println("ERROR: Empty File Name. Use `@` for using no File Name ...")
		return 989
	} //end if
	if(fName != "#") {
		fName = filepath.Base(strings.TrimSpace(fName))
		fName = strings.TrimSpace(fName)
		if(fName == "") {
			log.Println("ERROR: Invalid File Name")
			return 988
		}
	} //end if
	//--
	bar := pb.New(int(datalen)).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 250)
	bar.ShowSpeed = true
	bar.Start()
	//--
	var buffer bytes.Buffer
	w := multipart.NewWriter(&buffer)
	w.CreateFormField("test")
	if(fName != "#") {
		w.WriteField("test", "http post with file")
		fw, err := w.CreateFormFile("file", fName)
		if(err != nil) {
			log.Println("ERROR: Failed to Create Form Field: file")
			return 977
		} //end if
		_, err = io.Copy(fw, data)
		if(err != nil) {
			log.Println("ERROR: Failed to Populate Form Field: file")
			return 978
		} //end if
	} else {
		w.WriteField("test", "http post")
		_, err := w.CreateFormField("data")
		if(err != nil) {
			log.Println("ERROR: Failed to Create Form Field: data")
			return 977
		} //end if
		dt := new(bytes.Buffer)
		dt.ReadFrom(data)
		w.WriteField("data", dt.String())
	} //end if
	//--
	w.Close() // IMPORTANT: if you do not close the multipart writer you will not have a terminating boundry
	//--
	var resp *http.Response
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(http.MethodPost, url, bar.NewProxyReader(&buffer))
	req.Close = true
	if(err != nil) {
		log.Println("ERROR: Failed to handle HTTP Client: ", err)
		return 999
	} //end if
	req.Header.Set("Content-Type", w.FormDataContentType())
	//--
	if(username != "") {
		req.SetBasicAuth(username, passwd)
	} //end if
	//--
	resp, err = client.Do(req)
	if(err != nil) {
		log.Println("ERROR: Failed to handle HTTP Client Request: ", err)
		return 998
	} //end if
	resp.Body.Close()
	//--
//	bar.Finish()
	bar.FinishPrint("Data Form Post Completed: " + url)
	//--
	return resp.StatusCode
	//--
} //END FUNCTION


func testPostText() int {
	//-- sample POST text
	var txtStr string = "any thing"
	//--
	log.Println("Sample: HTTP POST Text: `" + txtStr + "` / Length:", len(txtStr), "bytes")
	//--
	return postRequest(THE_AUTH_USERNAME, THE_AUTH_PASSWORD, THE_URL, "#", strings.NewReader(txtStr), int64(len(txtStr)))
	//--
} //END FUNCTION


func testPostJson() int {
	//-- sample POST json
	var jsonStr string = `{"name":"Rob", "title":"developer"}`
	//--
	log.Println("Sample: HTTP POST Json: `" + jsonStr + "` / Length:", len(jsonStr), "bytes")
	//--
	return postRequest(THE_AUTH_USERNAME, THE_AUTH_PASSWORD, THE_URL, "#", bytes.NewBuffer([]byte(jsonStr)), int64(len(jsonStr)))
	//--
} //END FUNCTION


func testPostFile(fPath string) int {
	//--
	var fName string = strings.TrimSpace(fPath)
	if(fName == "") {
		log.Println("Empty File Name to Upload")
		os.Exit(1)
	} //end if
	fName = filepath.Base(fName)
	if(fName == "") {
		log.Println("Invalid File Name to Upload")
		os.Exit(1)
	} //end if
	//--
	theFile := "test/" + fName
	log.Println("Sample: HTTP POST File: `" + theFile + "`")
	//-- test if file exists, get the file size and read the file (open)
	fi, err := os.Stat(theFile)
	if(err != nil) {
		log.Fatal(err)
	} //end if
	var size int64 = fi.Size()
	data, err := os.Open(theFile)
	if(err != nil) {
		log.Fatal(err)
	} //end if
	defer data.Close()
	//--
	return postRequest(THE_AUTH_USERNAME, THE_AUTH_PASSWORD, THE_URL, theFile, data, size)
	//--
} //END FUNCTION


func main()  {
	//--
	if(len(os.Args) < 2) {
		log.Println("Usage: http-post-request text|json|file")
		os.Exit(1)
	} //end if
	//--
	var httpRequestStatusCode int = 0
	//--
	switch(os.Args[1]) {
		case "text":
			httpRequestStatusCode = testPostText()
			break
		case "json":
			httpRequestStatusCode = testPostJson()
			break
		case "file":
			if(len(os.Args) != 3) {
				log.Println("Usage: http-post-request file test/<filename>")
				os.Exit(1)
			} //end if
			httpRequestStatusCode = testPostFile(os.Args[2])
			break
		default:
			log.Println("Invalid mode. Invoke -help to see the arguments ...")
			os.Exit(2)
	} //end switch
	//--
	if(httpRequestStatusCode == 200) {
		fmt.Println("HTTP 200 OK")
	} else if(httpRequestStatusCode == 201) {
		fmt.Println("HTTP 201 CREATED")
	} else if(httpRequestStatusCode == 202) {
		fmt.Println("HTTP 202 ACCEPTED")
	} else if(httpRequestStatusCode == 400) {
		fmt.Println("HTTP 400 BAD REQUEST")
	} else if(httpRequestStatusCode == 401) {
		fmt.Println("HTTP 401 UNAUTHORIZED (Authentication Failed)")
	} else if(httpRequestStatusCode == 403) {
		fmt.Println("HTTP 403 FORBIDDEN (Access Denied)")
	} else {
		fmt.Println("HTTP Status NOT OK / Code: " + strconv.Itoa(httpRequestStatusCode))
	} //end if else
	//--
} //END FUNCTION

// #END
