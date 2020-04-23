
// GO Lang
// Markers TPL dev

package main

import (
	"os"
	"log"
	"fmt"
	smart "github.com/unix-world/smartgo"
	uid  "github.com/unix-world/smartgo/uuid"
)

const (
	THE_TPL = `Hallo, this is Markers TPL: [###MARKER|json###] [###MARKER2|url|html###]
`
)


func main() {

//                   1234567890_ abcdefghijklmnopqrstuvwxyz-ABCDEFGHIJKLMNOPQRSTUVWXYZ:;"'~`!@#$%^&*()+=[]{}|\<>,.?/

	var testStr = "1234567890_ abcdefghijklmnopqrstuvwxyz-ABCDEFGHIJKLMNOPQRSTUVWXYZ:;\"'~`!@#$%^&*()+=[]{}|\\<>,.?/\t\r\n@"
	testStr = " Lorem Ipsum șȘțȚâÂăĂîÎ is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.       \r\n      \r\n\r\n"
	testStr = `{"a":2, "b":"<c d=\"About:https://a.b?d=1&c=2\">"}`

	if(smart.StrGetUnicodeSubstring(testStr, 0, 0) != smart.StrGetAsciiSubstring(testStr, 0, 0)) {
		log.Fatal("SubString Comparison Test Failed !")
		os.Exit(1)
	}

	var arr = map[string]string{
		"MARKER": 	testStr, //"aA-șȘțȚâÂăĂîÎ_+100.12345678901",
		"MARKER2": 	`<Tyler="test">`,
	}

	var u = uid.Uuid()
	u = uid.Uuid()
	u = uid.Uuid()
	tpl := smart.RenderMarkersTpl(THE_TPL, arr, false, false)
	fmt.Println("UUID:", u, "\n" + "TPL: " + "\n" + tpl)

} //END FUNCTION

// #END
