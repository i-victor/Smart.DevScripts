
package main

import (
	"os"
	"log"
	"fmt"
	"time"
//	"bytes"
	"strings"
	"strconv"
	"io"
	"path/filepath"
	"net/http"
	"crypto/tls"

	"github.com/cheggaaa/pb"
)


func putRequest(username string, passwd string, url string, fName string, data io.Reader, datalen int64) int {
	//--
	fName = filepath.Base(strings.TrimSpace(fName))
	if(fName == "") {
		log.Println("ERROR: Empty File Name")
		return 999
	}
	fName = strings.TrimSpace(fName)
	if(fName == "") {
		log.Println("ERROR: Invalid File Name")
		return 998
	}
	//--
	bar := pb.New(int(datalen)).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 250)
	bar.ShowSpeed = true
	bar.Start()
	//--
	var realURL string = strings.TrimRight(url, "/") + "/" + fName
	var resp *http.Response
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(http.MethodPut, realURL,  bar.NewProxyReader(data))
	if err != nil {
		log.Println("ERROR: Failed to handle HTTP Client: ", err)
		return 997
	}
	if username != "" {
		req.SetBasicAuth(username, passwd)
	}
	req.TransferEncoding = []string{"identity"} // forces to change the default chunked transfer encoding and set it to gzip (support wider servers)
	req.ContentLength = datalen
	//--
	resp, err = client.Do(req)
	if err != nil {
		log.Println("ERROR: Failed to handle HTTP Client Request: ", err)
	}
	//--
//	bar.Finish()
	bar.FinishPrint("Data Transfer Completed: " + realURL)
	//--
	return resp.StatusCode
	//--
}


func main()  {

//	putRequest("http://127.0.0.1", strings.NewReader("any thing"))

//	var jsonStr string = `{"name":"Rob", "title":"developer"}`
//	putRequest("http://127.0.0.1/test-json-put.php", bytes.NewBuffer([]byte(jsonStr)))

	if len(os.Args) != 2 {
		log.Println("Usage: http-put-request test/<filename>")
		os.Exit(1)
	}

	var fName string = strings.TrimSpace(os.Args[1])
	if(fName == "") {
		log.Println("Empty File Name to Upload")
		os.Exit(1)
	}
	fName = filepath.Base(fName)
	if(fName == "") {
		log.Println("Invalid File Name to Upload")
		os.Exit(1)
	}

	theFile := "test/" + fName

	// test if file exists, get the file size and read the file (open)
	fi, err := os.Stat(theFile);
	if err != nil {
		log.Fatal(err)
	}
	var size int64 = fi.Size()
	data, err := os.Open(theFile)
	if err != nil {
		log.Fatal(err)
	}

//	var httpRequestStatusCode int = putRequest("admin", "pass", "http://d1.softlandro.com:13080/webdav/", theFile, data, size)
	var httpRequestStatusCode int = putRequest("admin", "pass", "https://d1.softlandro.com:13443/webdav/", theFile, data, size)

	if(httpRequestStatusCode == 200) {
		fmt.Println("HTTP 200 OK");
	} else if(httpRequestStatusCode == 201) {
		fmt.Println("HTTP 201 ACCEPTED");
	} else if(httpRequestStatusCode == 400) {
		fmt.Println("HTTP 400 BAD REQUEST");
	} else if(httpRequestStatusCode == 401) {
		fmt.Println("HTTP 401 UNAUTHORIZED (Authentication Failed)");
	} else if(httpRequestStatusCode == 403) {
		fmt.Println("HTTP 403 FORBIDDEN (Access Denied)");
	} else {
		fmt.Println("HTTP Status NOT OK / Code: " + strconv.Itoa(httpRequestStatusCode));
	} //end if else

}