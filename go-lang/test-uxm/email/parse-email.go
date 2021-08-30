
// GO Lang
// email / mime parser
// (c) 2017-2018 unix-world.org
// version: 2018.12.02

package main

import (
	"fmt"
	"log"
	"strings"

//	parsemail "github.com/DusanKasan/parsemail"
	parsemail "github.com/unix-world/smartgo/parsemail" // fixed version
	smart "github.com/unix-world/smartgo"
)


func main() {

	str, _ := smart.SafePathFileRead("./message3.eml", false)
	str = smart.StrTrimWhitespaces(str) // important: if there are new lines at the begining of mime message will fail to get content type

	var reader = strings.NewReader(str)
	email, err := parsemail.Parse(reader) // returns Email struct and error
	if err != nil {
		log.Println("ERROR:", err)
	}

	fmt.Println("ContentType:", email.ContentType)
	fmt.Println("Subject", email.Subject)
	fmt.Println("From", email.From)
	fmt.Println("To", email.To)
	fmt.Println("Body:", "\n", "`" + email.HTMLBody + "`")

}

// #END
