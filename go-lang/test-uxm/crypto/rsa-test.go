
package main

import (
	"os"
	"log"

	"crypto"
	"crypto/rand"
	"crypto/sha512"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

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

	LogToConsoleWithColors()

	// The GenerateKey method takes in a reader that returns random bits, and
	// the number of bits
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096) // 4096 bit = 512 bytes
	if err != nil {
		log.Println("[ERROR] Cannot generate RSA Private Key:", err)
		return
	}

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if(err != nil) {
		log.Println("[ERROR] Cannot marshal RSA Private Key:", err)
		return
	}
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})
	log.Println("[DATA] Private PEM:", "\n" + string(privPEM), smart.StrChunkSplit(smart.Base64Encode(string(privBytes)), 64, "\n"))

	pubBytes, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if(err != nil) {
		log.Println("[ERROR] Cannot marshal RSA Public Key:", err)
		return
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})
	log.Println("[DATA] Public PEM:", "\n" + string(pubPEM), smart.StrChunkSplit(smart.Base64Encode(string(pubBytes)), 64, "\n"))

	// The public key is a part of the *rsa.PrivateKey struct
	publicKey := privateKey.PublicKey

	//====================
	// RSA is only able to encrypt data to a maximum amount equal to your key size
	// 1024 bits = 128 bytes
	// 2048 bits = 256 bytes
	// 4096 bits = 512 bytes
	// from the above values have to substract padding/header data (see below)
	//====================
	// RSA/ECB/NoPadding, 0
	// RSA/ECB/PKCS1Padding, 11
	// RSA/ECB/OAEPPadding, 42 // Actually it's OAEPWithSHA1AndMGF1Padding
	// RSA/ECB/OAEPWithMD5AndMGF1Padding, 34
	// RSA/ECB/OAEPWithSHA1AndMGF1Padding, 42
	// RSA/ECB/OAEPWithSHA224AndMGF1Padding, 58
	// RSA/ECB/OAEPWithSHA256AndMGF1Padding, 66
	// RSA/ECB/OAEPWithSHA384AndMGF1Padding, 98
	// RSA/ECB/OAEPWithSHA512AndMGF1Padding, 130
	// RSA/ECB/OAEPWithSHA3-224AndMGF1Padding, 58
	// RSA/ECB/OAEPWithSHA3-256AndMGF1Padding, 66
	// RSA/ECB/OAEPWithSHA3-384AndMGF1Padding, 98
	// RSA/ECB/OAEPWithSHA3-512AndMGF1Padding, 130
	//====================
	var sMsg string = smart.StrSubstr(strings.Repeat("super secret message" + "\n", 100), 0, 512-130)
	encryptedBytes, err := rsa.EncryptOAEP(
		sha512.New(),
		rand.Reader,
		&publicKey,
		[]byte(sMsg),
		nil)
	if err != nil {
		fatalError(err)
		return
	}

	log.Println("[DATA] encrypted message:", smart.Base64Encode(string(encryptedBytes)))

	// The first argument is an optional random data generator (the rand.Reader we used before)
	// we can set this value as nil
	// The OAEPOptions in the end signify that we encrypted the data using OAEP, and that we used
	// SHA512 to hash the input.
	decryptedBytes, err := privateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA512})
	if err != nil {
		fatalError(err)
		return
	}

	// We get back the original information in the form of bytes, which we
	// the cast to a string and print
	log.Println("[DATA] decrypted message: ", string(decryptedBytes))

	if(string(decryptedBytes) != sMsg) {
		log.Println("[ERROR] Keys are different")
		return
	}
	log.Println("[OK] Keys are similar")


/*
msg := []byte("verifiable message")

// Before signing, we need to hash our message
// The hash is what we actually sign
msgHash := sha512.New()
_, err = msgHash.Write(msg)
if err != nil {
	panic(err)
}
msgHashSum := msgHash.Sum(nil)

// In order to generate the signature, we provide a random number generator,
// our private key, the hashing algorithm that we used, and the hash sum
// of our message
signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA512, msgHashSum, nil)
if err != nil {
	panic(err)
}

// To verify the signature, we provide the public key, the hashing algorithm
// the hash sum of our message and the signature we generated previously
// there is an optional "options" parameter which can omit for now
err = rsa.VerifyPSS(&publicKey, crypto.SHA512, msgHashSum, signature, nil)
if err != nil {
	fmt.Println("could not verify signature: ", err)
	return
}
// If we don't get any error from the `VerifyPSS` method, that means our
// signature is valid
fmt.Println("signature verified")
*/

}

// #end
