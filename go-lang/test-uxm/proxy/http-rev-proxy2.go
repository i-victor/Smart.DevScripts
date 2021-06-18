
package main

import(
	"log"
	"fmt"
	"net/url"
	"net/http"
	"net/http/httputil"
	"io/ioutil"
	"bytes"
	"strconv"
)

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
	//--
	// IMPORTANT: the body is generally available on GET or POST ; by example on HEAD the body is empty, aka 304
	//-- modify body
//	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
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
	remote, err := url.Parse("http://sf.loc:80/")
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Transport = &transport{http.DefaultTransport}

	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		r.Host = r.URL.Host // ! important ! to send the request to the correct host
		w.Header().Set("X-Real-IP", "127.0.0.1")
		p.ServeHTTP(w, r)
	}
}

// END
