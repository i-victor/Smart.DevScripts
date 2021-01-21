
// GoLang Sample
// HTTP Client PUT : stream file | json | text
// (c) 2020-2021 unix-world.org
// r.20210118.2155

package main

import (
	"os"
	"log"
	"fmt"
	"time"
	"bytes"
	"strings"
	"strconv"
	"io"
	"path/filepath"
	"net/http"
	"crypto/tls"

	"github.com/cheggaaa/pb"
)


const ( // all tests can be performed with: webdav-server.go which serves WebDAV protocol (on /webdav/) or plain HTTP(S) (on /)
	THE_URL = "http://127.0.0.1:13080/" 		// Ex(text|json): `http://127.0.0.1:13080/` or `https://127.0.0.1:13443/` ; Ex(file): `http://127.0.0.1:13080/webdav/` or `https://127.0.0.1:13443/webdav/`
	THE_AUTH_USERNAME = "admin"					// leave empty if no auth required or fill the auth username ; for webdav the default is `admin`
	THE_AUTH_PASSWORD = "pass"					// leave empty if no auth required or fill the auth password ; for webdav the default is `pass`
)


func putRequest(username string, passwd string, url string, fName string, data io.Reader, datalen int64) int {
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
	var realURL string = url
	if(fName != "#") {
		realURL = strings.TrimRight(realURL, "/") + "/" + fName
	} //end if
	//--
	var resp *http.Response
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(http.MethodPut, realURL,  bar.NewProxyReader(data))
	req.Close = true
	if(err != nil) {
		log.Println("ERROR: Failed to handle HTTP Client: ", err)
		return 999
	} //end if
	//--
	if(username != "") {
		req.SetBasicAuth(username, passwd)
	} //end if
	req.TransferEncoding = []string{"identity"} // forces to change the default chunked transfer encoding and set it to gzip (support wider servers)
	req.ContentLength = datalen
	//--
	resp, err = client.Do(req)
	if(err != nil) {
		log.Println("ERROR: Failed to handle HTTP Client Request: ", err)
		return 998
	} //end if
	resp.Body.Close()
	//--
//	bar.Finish()
	bar.FinishPrint("Data Transfer Completed: " + realURL)
	//--
	return resp.StatusCode
	//--
} //END FUNCTION


func testPutText() int {
	//-- sample PUT text
	var txtStr string = "any thing"
	//--
	log.Println("Sample: HTTP PUT Text: `" + txtStr + "` / Length:", len(txtStr), "bytes")
	//--
	return putRequest(THE_AUTH_USERNAME, THE_AUTH_PASSWORD, THE_URL, "#", strings.NewReader(txtStr), int64(len(txtStr)))
	//--
} //END FUNCTION


func testPutJson() int {
	//-- sample PUT json
	var jsonStr string = `{"name":"Rob", "title":"developer"}`
	//--
	log.Println("Sample: HTTP PUT Json: `" + jsonStr + "` / Length:", len(jsonStr), "bytes")
	//--
	return putRequest(THE_AUTH_USERNAME, THE_AUTH_PASSWORD, THE_URL, "#", bytes.NewBuffer([]byte(jsonStr)), int64(len(jsonStr)))
	//--
} //END FUNCTION


func testPutFile(fPath string) int {
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
	log.Println("Sample: HTTP PUT File: `" + theFile + "`")
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
	return putRequest(THE_AUTH_USERNAME, THE_AUTH_PASSWORD, strings.TrimRight(THE_URL, "/") + "/webdav/", theFile, data, size)
	//--
} //END FUNCTION


func main()  {
	//--
	if(len(os.Args) < 2) {
		log.Println("Usage: http-put-request text|json|file")
		os.Exit(1)
	} //end if
	//--
	var httpRequestStatusCode int = 0
	//--
	switch(os.Args[1]) {
		case "text":
			httpRequestStatusCode = testPutText()
			break
		case "json":
			httpRequestStatusCode = testPutJson()
			break
		case "file":
			if(len(os.Args) != 3) {
				log.Println("Usage: http-put-request file test/<filename>")
				os.Exit(1)
			} //end if
			httpRequestStatusCode = testPutFile(os.Args[2])
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
