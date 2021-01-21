
// GoLang Sample
// HTTP Client GET
// (c) 2020-2021 unix-world.org
// r.20210118.2155

package main

import (
//	"os"
	"log"
	"fmt"
	"strconv"
	"net/http"
	"crypto/tls"
)


const ( // all tests can be performed with: webdav-server.go which also serves plain HTTP(S) (on /)
	THE_URL = "http://127.0.0.1:8087/" 			// `http://127.0.0.1:80/` or `https://127.0.0.1:443/`
	THE_AUTH_USERNAME = ""						// leave empty if no auth required or fill the auth username
	THE_AUTH_PASSWORD = ""						// leave empty if no auth required or fill the auth password
)


func getRequest(username string, passwd string, url string) int {
	//--
	var resp *http.Response
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Close = true
	if(err != nil) {
		log.Println("ERROR: Failed to handle HTTP Client: ", err)
		return 999
	} //end if
	req.Header.Add("x-test", "This is a test from GoLang ...")
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
	return resp.StatusCode
	//--
} //END FUNCTION


func testGet() int {
	//--
	log.Println("Sample: HTTP GET")
	//--
	return getRequest(THE_AUTH_USERNAME, THE_AUTH_PASSWORD, THE_URL)
	//--
} //END FUNCTION


func main()  {
	//--
	var httpRequestStatusCode int = testGet()
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
