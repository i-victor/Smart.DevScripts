package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"

	smart "github.com/unix-world/smartgo"
)

func main() {

	str, _ := smart.SafePathFileRead("./message2.eml", false)
	str = smart.StrTrimWhitespaces(str) // important: if there are new lines at the begining of mime message will fail to get content type

	msg, err := mail.ReadMessage(bytes.NewBufferString(str))
	if err != nil {
		log.Fatal("Cannot parse Message.")
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(mediaType)

	from, _ := (&mail.AddressParser{}).ParseList(msg.Header.Get("From"))
	fmt.Println("From", from)

	to, _ := (&mail.AddressParser{}).ParseList(msg.Header.Get("To"))
	fmt.Println("To", to)

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(msg.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatal(err)
			}

			// decode any Q-encoded values
			for name, values := range p.Header {
				for idx, val := range values {
					fmt.Printf("%d: %s: %s\n", idx, name, decodeRFC2047(val))
				}
			}

			slurp, err := ioutil.ReadAll(p)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Content-type: %s\n%s\n", p.Header.Get("Content-Type"), slurp)
		}
	}
}

// decodeRFC2047 ...
func decodeRFC2047(s string) string {
	// GO 1.5 does not decode headers, but this may change in future releases...
	decoded, err := (&mime.WordDecoder{}).DecodeHeader(s)
	if err != nil || len(decoded) == 0 {
		return s
	}
	return decoded
}
