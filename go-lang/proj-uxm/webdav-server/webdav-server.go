
// small webdav server written in go

package main

import (
	"os"
	"log"
	"fmt"
	"flag"
	"strconv"
	"bytes"
	"encoding/base64"
	"html"
	"net/http"
	"crypto/subtle"

	"golang.org/x/net/webdav"
)
const
(
	THE_VERSION = "r.20200421.2337"
	STORAGE_DIR = "./dav-storage"
	DAV_PATH = "/webdav"
	CONN_HOST = "0.0.0.0"
	CONN_PORT = 13080
	CONN_SPORT = 13443
	ADMIN_USER = "admin"
	ADMIN_PASSWORD = "pass"
	SVG_SPIN = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="32" height="32" fill="grey" id="loading-spin-svg"><path opacity=".25" d="M16 0 A16 16 0 0 0 16 32 A16 16 0 0 0 16 0 M16 4 A12 12 0 0 1 16 28 A12 12 0 0 1 16 4"/><path d="M16 0 A16 16 0 0 1 32 16 L28 16 A12 12 0 0 0 16 4z"><animateTransform attributeName="transform" type="rotate" from="0 16 16" to="360 16 16" dur="0.8s" repeatCount="indefinite" /></path></svg>`
)

func main() {

	//--

	path, err := os.Getwd()
	if err != nil {
		log.Println("ERROR: Cannot Get Current Path: ", err)
		return
	}

	//--

	httpAddr := flag.String("h", CONN_HOST, "HTTP Host to listen to. The 0.0.0.0 will listen on all interfaces. 127.0.0.1 will listen only on localhost. can be any of these or a local or internet valid IPv4")
	httpPort := flag.Int("p", CONN_PORT, "HTTP Port to listen to")
	httpsPort := flag.Int("ps", CONN_SPORT, "HTTPS Port to listen to in secure / TLS mode")
	serveSecure := flag.Bool("s", false, "Serve HTTPS. Default false. If False will serve just HTTP")
	disableUnsecure := flag.Bool("u", false, "Disable Serve HTTP and Serve only HTTPS. Default false. If True will serve just HTTPS")

	flag.Parse()

	//-- for web

	var serverSignature bytes.Buffer
	serverSignature.WriteString("MiniWebDAV GO Server " + THE_VERSION + "\n")
	serverSignature.WriteString("(c) 2020 unix-world.org" + "\n")
	serverSignature.WriteString("\n")
	if *disableUnsecure != true {
		serverSignature.WriteString("<URL> :: http://" + *httpAddr + ":" + strconv.Itoa(*httpPort) + DAV_PATH + "/" + "\n")
	}
	if *serveSecure == true {
		serverSignature.WriteString("<Secure URL> :: https://" + *httpAddr + ":" + strconv.Itoa(*httpsPort) + DAV_PATH + "/" + "\n")
	}

	//-- for console

	fmt.Println("===========================================================================")
	fmt.Println("MiniWebDAV GO Server " + THE_VERSION)
	fmt.Println("---------------------------------------------------------------------------")
	fmt.Println("Current Path: " + string(path))
	fmt.Println("DAV Folder: " + STORAGE_DIR)
	fmt.Println("---------------------------------------------------------------------------")
	if *disableUnsecure != true {
		fmt.Println("Listening at http://" + *httpAddr + ":" + strconv.Itoa(*httpPort) + "/")
	}
	if *serveSecure == true {
		fmt.Println("Secure TLS Listening at https://" + *httpAddr + ":" + strconv.Itoa(*httpsPort) + "/")
	}
	fmt.Println("===========================================================================")

	//--

	srv := &webdav.Handler{
		Prefix:     DAV_PATH,
		FileSystem: webdav.Dir(STORAGE_DIR),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("MiniWebDAV GO Server :: WEBDAV.ERROR: %s [%s %s %s] %s [%s] %s\n", err, r.Method, r.URL, r.Proto, "*", r.Host, r.RemoteAddr)
			} else {
				log.Printf("MiniWebDAV GO Server :: WEBDAV [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, "*", r.Host, r.RemoteAddr)
			}
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var statusCode = 202
		if r.URL.Path != "/" {
			statusCode = 404
			w.WriteHeader(statusCode)
			w.Write([]byte("404 Not Found\n"))
			log.Printf("MiniWebDAV GO Server :: DEFAULT.ERROR [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, strconv.Itoa(statusCode), r.Host, r.RemoteAddr)
			return
		}
		log.Printf("MiniWebDAV GO Server :: DEFAULT [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, strconv.Itoa(statusCode), r.Host, r.RemoteAddr)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(statusCode) // status code must be after content type
		fmt.Fprintf(w, "<!DOCTYPE html>" + "\n" + "<html>" + "\n" + "<head>" + "\n" + `<meta charset="UTF-8">` + "\n" + "<title>" + html.EscapeString("MiniWebDAV GO Server " + THE_VERSION) + "</title></head>" + "\n" + "<body>" + "\n" + `<div style="text-align:center; margin:10px; cursor:help;"><img alt="Status: Up and Running ..." title="Status: Up and Running ..." width="96" height="96" src="data:image/svg+xml;base64,` + base64.StdEncoding.EncodeToString([]byte(SVG_SPIN)) + `"></div>` + "\n" + `<div style="background:#778899; color:#FFFFFF; font-size:2rem; font-weight:bold; text-align:center; border-radius:3px; padding:10px; margin:20px;">` + "\n" + "<pre>" + "\n" + html.EscapeString(serverSignature.String()) + "</pre>" + "\n" + "</div>" + "\n" + "</body>" + "\n" + "</html>" + "\n")
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		var statusCode = 203
		log.Printf("MiniWebDAV GO Server :: VERSION [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, strconv.Itoa(statusCode), r.Host, r.RemoteAddr)
		// plain/text
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, "MiniWebDAV GO Server " + THE_VERSION + "\n")
	})

	//http.Handle(DAV_PATH+"/", srv)
	http.HandleFunc(DAV_PATH+"/", func(w http.ResponseWriter, r *http.Request) {
		// test if basic auth
		user, pass, ok := r.BasicAuth()
		// check if basic auth and if credentials match
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(ADMIN_USER)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(ADMIN_PASSWORD)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="MiniWebDAV GO Server Storage Area"`)
			w.WriteHeader(401) // status code must be after set headers
			w.Write([]byte("401 Unauthorized\n"))
			log.Printf("MiniWebDAV GO Server :: AUTH.FAILED [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, "401", r.Host, r.RemoteAddr)
			return
		}
		// if all ok above (basic auth + credentials ok, serve ...)
		srv.ServeHTTP(w, r)

	})

	if *serveSecure == true {
		if _, err := os.Stat("./pem/cert.pem"); err != nil {
			fmt.Println("[INIT TLS ERROR] No cert.pem in current directory. Please provide a valid cert")
			return
		}
		if _, er := os.Stat("./pem/key.pem"); er != nil {
			fmt.Println("[INIT TLS ERROR] No key.pem in current directory. Please provide a valid cert")
			return
		}
		if *disableUnsecure != true {
			fmt.Println("... starting TLS / HTTPS")
			go http.ListenAndServeTLS(*httpAddr + fmt.Sprintf(":%d", *httpsPort), "./pem/cert.pem", "./pem/key.pem", nil)
		} else {
			fmt.Println("... starting *ONLY* TLS / HTTPS")
			http.ListenAndServeTLS(*httpAddr + fmt.Sprintf(":%d", *httpsPort), "./pem/cert.pem", "./pem/key.pem", nil)
		}
	}
	if *disableUnsecure != true {
		fmt.Println("... starting HTTP")
		if err := http.ListenAndServe(*httpAddr + fmt.Sprintf(":%d", *httpPort), nil); err != nil {
			log.Fatalf("ERROR: WebDAV Server: %v", err)
		}
	}
	if *disableUnsecure == true && *serveSecure != true {
		fmt.Println("[INIT ERROR] The both HTTP and HTTPS modes are disabled ... server will exit ...")
	}

}


// #end
