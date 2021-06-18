package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"crypto/tls"
	"io/ioutil"
	"bytes"
	"strconv"
)

// # https://www.integralist.co.uk/posts/golang-reverse-proxy/

// by unixman

type transport struct {
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	//-- modify body
//	b = bytes.Replace(b, []byte("Core Test Suite"), []byte("Core TEST Suite"), -1)
	//--
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	//--
	fmt.Println("====================================================", "Status:", resp.Status)
	fmt.Println(string(b)) // print the body to stdout
	fmt.Println("====================================================")
	//--
	return resp, nil
}

func main() {

//	var theURL string = "https://127.0.0.1:443/"
	var theURL string = "http://sf.loc:80/"

	target, _ := url.Parse(theURL)

	fmt.Println("Schema:", target.Scheme)
	fmt.Println("Host:", target.Host)

//	proxy := httputil.NewSingleHostReverseProxy(target) // this does not rewrite the Host header ; To rewrite Host headers, use ReverseProxy directly with a custom Director policy as below

	director := func(req *http.Request) {

		rAddr, _, _ := net.SplitHostPort(req.RemoteAddr) // get only IP from IP:Port returned by req.RemoteAddr
	//	fmt.Println(rAddr)

		// see: https://en.wikipedia.org/wiki/List_of_HTTP_header_fields

		//-- originals
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		//--
		req.Header.Add("X-Origin-Host", target.Host)
		req.Header.Add("X-Forwarded-Host", req.Host)
		//--

	//	req.Header.Add("X-Forwarded-Host", req.URL.Host)
	//	req.Header.Add("X-Origin-Host", target.Host)
	//	req.Header.Add("X-Forwarded-Host", target.Host)

	//	req.Header.Add("Client-IP", "10.0.0.2")
	//	req.Header.Add("X-Client-IP", "10.0.0.1")
		req.Header.Add("X-Forwarded-For", rAddr)
		req.Header.Add("X-Real-IP", "127.0.0.1")

	}

//	proxy := &httputil.ReverseProxy{Director: director, Transport:&transport{http.DefaultTransport}}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.Transport = &transport{http.DefaultTransport}

	//-- TLS flexible, allow insecure
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Host = r.URL.Host // ! important ! to send the request to the correct host
		proxy.ServeHTTP(w, r)
	})

//	log.Fatal(http.ListenAndServe(":8888", nil))
	log.Fatal(http.ListenAndServeTLS(":8888", "server.crt", "server.key", nil))

}

