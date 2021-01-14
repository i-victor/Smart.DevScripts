
// GoLang
// HTTP Client GET : Hit Timer
// (c) 2020-2021 unix-world.org
// r.20210114.1634

package main

import (
//	"os"
	"log"
	"fmt"
	"strconv"
	"net/http"
	"crypto/tls"
	"time"
)


const ( // all tests can be performed with: webdav-server.go which also serves plain HTTP(S) (on /)
	THE_URL = "http://127.0.0.1:80/" 			// `http://127.0.0.1:80/` or `https://127.0.0.1:443/`
	THE_AUTH_USERNAME = ""						// leave empty if no auth required or fill the auth username
	THE_AUTH_PASSWORD = ""						// leave empty if no auth required or fill the auth password

	WAIT_STATUS_CODE_HTTP = 429 				// 429 is too many requests, flag to wait
	INTERVAL_HIT = 1							// interval in seconds to wait between hit requests
	INTERVAL_WAIT = 3 							// interval in seconds to wait if the WAIT_STATUS_CODE_HTTP is get by hit
)


func getRequest(username string, passwd string, url string) int {
	//--
	var resp *http.Response
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if(err != nil) {
		log.Println("ERROR: Failed to handle HTTP Client: ", err)
		return 997
	} //end if
	//--
	if(username != "") {
		req.SetBasicAuth(username, passwd)
	} //end if
	//--
	resp, err = client.Do(req)
	if(err != nil) {
		log.Println("ERROR: Failed to handle HTTP Client Request: ", err)
	} //end if
	//--
	return resp.StatusCode
	//--
} //END FUNCTION


func hitTheUrl() int {
	//--
	log.Println("HTTP GET", THE_URL)
	if(THE_AUTH_USERNAME != "") {
		log.Println("HTTP AUTH", "User:`" + THE_AUTH_USERNAME, "` ; Pass:[" + strconv.Itoa(len(THE_AUTH_PASSWORD)) + "]")
	} //end if
	//--
	var httpRequestStatusCode int = getRequest(THE_AUTH_USERNAME, THE_AUTH_PASSWORD, THE_URL)
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
	return httpRequestStatusCode
	//--
} //END FUNCTION


func main()  {
	//--
	var httpRequestStatusCode int = 0
	var hitIntervalSleep int = 0
	var theWaitInterval string = ""
	//--
	for i := 0; i >= 0; i++ { // infinite loop
		httpRequestStatusCode = hitTheUrl()
		fmt.Println("Hit # " + strconv.Itoa(i) + " :: HTTP Status Code: " + strconv.Itoa(httpRequestStatusCode))
		if(httpRequestStatusCode == WAIT_STATUS_CODE_HTTP) {
			hitIntervalSleep = INTERVAL_WAIT
			theWaitInterval = "Wait"
		} else {
			hitIntervalSleep = INTERVAL_HIT
			theWaitInterval = "Default"
		}
		if(hitIntervalSleep > 0) {
			time.Sleep(time.Duration(hitIntervalSleep) * time.Second)
			fmt.Println("@", theWaitInterval, "Interval (seconds):", hitIntervalSleep)
		}
		fmt.Println("========== [Done]", "\n\n")
	} //end for
	//--
} //END FUNCTION

// #END
