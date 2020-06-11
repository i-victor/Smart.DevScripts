
// Go Lang
// generate hexkey for ifconfig wifi in OpenBSD
// (c) 2020 unix-world.org

package main


import (
	"os"
	"fmt"
	"flag"
	"strings"
	"encoding/hex"
	"crypto/sha1"
	"golang.org/x/crypto/pbkdf2"
)


func main() {

	//--
	nwid := flag.String("nwid", "", "The wifi network ID")
	pass := flag.String("pass", "", "The wifi network password")
	//--
	flag.Parse()
	//--

	//--
	var theNwid string = *nwid
	var thePass string = *pass
	//--

	//--
	if((theNwid == "") || (thePass == "")) {
		fmt.Println("The Pass or Nwid must not be empty. See -help for more details ...")
		os.Exit(1)
	} //end if
	//--

	//--
	rawKey := pbkdf2.Key([]byte(thePass), []byte(theNwid), 4096, 32, sha1.New)
	//--
	var theHexKey string = strings.ToLower(hex.EncodeToString(rawKey))
	//--

	//--
	fmt.Println(`ifconfig iwm0 nwid "` + theNwid + `" wpakey 0x` + theHexKey)
	//--

} //END FUNCTION


// #END
