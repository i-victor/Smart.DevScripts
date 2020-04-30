
package main

import (
	"net/http"
	"crypto/tls"
	"strconv"
	"fmt"
	"io"
	"log"
//	"strings"
	"os"
//	"bytes"
)


func putRequest(url string, data io.Reader, datalen int64)  {
	var resp *http.Response
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
//	req.TransferEncoding = []string{"identity"} // changes the default chunked transfer encoding and set it to gzip
	req.ContentLength = datalen
	resp, err = client.Do(req)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	if(resp.StatusCode == 200) {
		fmt.Println("HTTP 200 OK");
	} else if(resp.StatusCode == 201) {
		fmt.Println("HTTP 201 ACCEPTED");
	} else {
		fmt.Println("HTTP Status NOT OK / Code: " + strconv.Itoa(resp.StatusCode));
	} //end if else
}


func main()  {

//	putRequest("http://127.0.0.1", strings.NewReader("any thing"))

//	var jsonStr string = `{"name":"Rob", "title":"developer"}`
//	putRequest("http://127.0.0.1/test-json-put.php", bytes.NewBuffer([]byte(jsonStr)))

	// read the file
	theFile := "test/page.pdf"
	fi, err := os.Stat(theFile);
	if err != nil {
		log.Fatal(err)
	}
	var size int64 = fi.Size()
	data, err := os.Open(theFile)
	if err != nil {
		log.Fatal(err)
	}
	putRequest("https://admin:the-pass@127.0.0.1/sites/smart-framework/admin.php/page/cloud.files/~/uploads/page.pdf", data, size)

}