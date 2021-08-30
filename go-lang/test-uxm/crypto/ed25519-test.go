
// GO Lang

package main

import (

	"os"
	"log"

	"crypto/rand"
	"crypto/ed25519"

	"strings"

	smart "github.com/unix-world/smartgo"
)

func LogToConsoleWithColors() {
	//--
//	smart.ClearPrintTerminal()
	//--
//	smart.LogToStdErr("DEBUG")
	smart.LogToConsole("DEBUG", true) // colored or not
//	smart.LogToFile("DEBUG", "logs/", "json", true, true) // json | plain ; also on console ; colored or not
	//--
//	log.Println("[DEBUG] Debugging")
//	log.Println("[DATA] Data")
//	log.Println("[NOTICE] Notice")
//	log.Println("[WARNING] Warning")
//	log.Println("[ERROR] Error")
//	log.Println("[OK] OK")
//	log.Println("A log message, with no type")
	//--
} //END FUNCTION


func fatalError(logMessages ...interface{}) {
	//--
	log.Fatal("[ERROR] ", logMessages) // standard logger
	os.Exit(1)
	//--
} //END FUNCTION


func main() {

	defer smart.PanicHandler()

	 LogToConsoleWithColors()

	var password = "a Password to encrypt the private ..."
	var encoding = "b92" // b92, b64s, b64, b62, b58, b36, b32, hex (b16)
	log.Println("[Info] Encoding:", encoding)

	pubKey, secKey, theErr := ed25519.GenerateKey(rand.Reader)
	if(theErr != nil) {
		log.Println("[ERROR]", theErr)
		return
	}
	var PublicKey string = smart.BaseEncode(pubKey, encoding)
	var SecretKey string = smart.ThreefishEncryptCBC(string(secKey), password, true)
	log.Println("[DATA] Public Key:", PublicKey)
	log.Println("[DATA] Private Key:", SecretKey)

	var pub2Key []byte = smart.BaseDecode(PublicKey, encoding)
	var sec2Key []byte = []byte(smart.ThreefishDecryptCBC(SecretKey, password, true))

	msg, errRd := smart.SafePathFileRead("ed25519-test.go", false)
	if(errRd != "") {
		log.Println("[ERROR]", errRd)
		return
	}
	if(msg == "") {
		log.Println("[ERROR]", "Message to sign is empty !")
		return
	}
	msg = strings.Repeat(msg, 1000)

	mLen := len(msg)
	log.Println("[INFO]: Message Length (bytes)", mLen)

	sigData := ed25519.Sign(sec2Key[:], []byte(msg))
	var SignatureData string = smart.BaseEncode(sigData, encoding)
	log.Println("[DATA] Signature Data:", SignatureData)
	var sig2Data []byte = smart.BaseDecode(SignatureData, encoding)

	sigSignature := ed25519.Sign(sec2Key[:], []byte(SignatureData))
	var SignatureOfSignature string = smart.BaseEncode(sigSignature, encoding)
	log.Println("[DATA] Signature Of Signature:", SignatureOfSignature)
	var sig2Signature []byte = smart.BaseDecode(SignatureOfSignature, encoding)

	if(ed25519.Verify(pub2Key[:], []byte(SignatureData), sig2Signature[:]) != true) {
		log.Println("[ERROR]: Signature Of Signature does not match")
		return
	} //end if
	log.Println("[OK]: Signature Of Signature match")

//msg = msg + " "
	if(ed25519.Verify(pub2Key[:], []byte(msg), sig2Data[:]) != true) {
		log.Println("[ERROR]: Signature does not match")
		return
	} //end if
	log.Println("[OK]: Signature match")

}


// #END
