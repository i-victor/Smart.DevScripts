
// GO Lang
// Markers TPL dev
// r.20200425.1525 :: STABLE

package main

import (
	"os"
	"log"
	"fmt"
	smart "github.com/unix-world/smartgo"
	uid  "github.com/unix-world/smartgo/uuid"
)

const (
	THE_TPL = `[%%%COMMENT%%%]
This is comment one ...
[%%%/COMMENT%%%]
Hallo, this is Markers TPL:[%%%|N%%%][###MARKER|json###][%%%|T%%%][###MARKER2|url|html###]
[%%%COMMENT%%%]
This is another comment ...
[%%%/COMMENT%%%]`
)


func main() {

	smart.LogToConsoleWithColors();

	input := "Lorem Ipsum dolor sit Amet"

	fmt.Println("MD5:", smart.Md5(input))
	fmt.Println("SHA1:", smart.Sha1(input))
	fmt.Println("SHA256:", smart.Sha256(input))
	fmt.Println("SHA384:", smart.Sha384(input))
	fmt.Println("SHA512:", smart.Sha512(input))

	b64 := smart.Base64Encode(input)
	fmt.Println("B64-Enc:", b64)
	fmt.Println("B64-Dec:", smart.Base64Decode(b64))

	hex := smart.Bin2Hex(input)
	fmt.Println("HEX-Enc:", hex)
	fmt.Println("HEX-Dec:", smart.Hex2Bin(hex))

	arch := smart.DataArchive(input)
	fmt.Println("Data-Arch:", arch)
	fmt.Println("Data-UnArch:", smart.DataUnArchive(arch))

	// INFO: arch data difers a little from PHP, maybe by some zlib metadata, but decrypt must work
	testPhpArchData := `HclBDkBAEETRw1hLplupZimDSMRKHMD06Psfgdj9/IfM1ZQ9Z00YLVlnfxNc+Zt+j6Phc+HM3tDkbcn7eR3tuU3SDKGhjwrCUaM4i6dbS7r9qRgEdIsq6i8=` + "\n" + `PHP.SF.151129/B64.ZLibRaw.HEX`;
	fmt.Println("Data-Arch-PHP:", testPhpArchData)
	testPhpUnArchData := smart.DataUnArchive(testPhpArchData)
	fmt.Println("Data-UnArch-PHP:", testPhpUnArchData)
	if(testPhpUnArchData != input) {
		log.Println("ERROR: DataArchive TEST Failed ... Archived Data is NOT EQUAL with Archived Data from PHP")
	} //end if

	//--
	bfKey := "some.BlowFish! - Key@Test 2ks i782s982 s2hwgsjh2wsvng2wfs2w78s528 srt&^ # *&^&#*# e3hsfejwsfjh"
	//--
	bfInput := input + " " + smart.DateNowUtc()
	fmt.Println("Data-To-Encrypt:", bfInput)
	bfData := smart.BlowfishEncryptCBC(bfInput, bfKey)
	fmt.Println("Data-Encrypted:", bfData)
	testDecBfData := smart.BlowfishDecryptCBC(bfData, bfKey)
	fmt.Println("Data-Decrypted:", testDecBfData)
	if((testDecBfData != bfInput) || (smart.Sha1(testDecBfData) != smart.Sha1(bfInput))) {
		log.Println("ERROR: BlowfishEncryptCBC TEST Failed ... Decrypted Data is NOT EQUAL with Plain Data")
	} //end if
	//--
	testPhpBfData := `695C491EF3E92DD8975423A91460F05F9DBBFDBE91DC55AE1D96CC43747B096D64CE08F42885D792505A56DF40CEE6B51FC399A3D756FADB4CE9A492BAE157B4B0DB0C6197D0E35B4C69F99266965686CB41628B75EA56CE006518F408CC0AF1`
	if(smart.BlowfishEncryptCBC(input, bfKey) != testPhpBfData) {
		log.Println("ERROR: BlowfishEncryptCBC TEST Failed ... Encrypted Data is NOT EQUAL with Encrypted Data from PHP")
	} //end if
	fmt.Println("Data-Encrypted-PHP:", testPhpBfData)
	testDecPhpBfData := smart.BlowfishDecryptCBC(testPhpBfData, bfKey)
	fmt.Println("Data-Decrypted-PHP:", testDecPhpBfData)
	if((testDecPhpBfData != input) || (smart.Sha1(testDecPhpBfData) != smart.Sha1(input))) {
		log.Println("ERROR: BlowfishDecryptCBC TEST Failed ... Decrypted Data is NOT EQUAL with Decrypted Data from PHP")
	} //end if


	//                   1234567890_ abcdefghijklmnopqrstuvwxyz-ABCDEFGHIJKLMNOPQRSTUVWXYZ:;"'~`!@#$%^&*()+=[]{}|\<>,.?/

	//-----

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
