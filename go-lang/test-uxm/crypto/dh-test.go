
package main

import (
	"os"
	"log"

	"github.com/monnand/dhkx"

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


//=======
// This is an implementation of Diffie-Hellman Key Exchange algorithm.
// The algorithm is used to establish a shared key between two communication peers without sharing secrete information.
// Typical process:
// First, Alice and Bob should agree on which group to use.
// If you are not sure, choose group 14. GetGroup() will return the desired group by a given id.
// GetGroup(0) will return a default group, which is usually safe enough to use this group.
// It is totally safe to share the group's information.
//=======

//const GroupID = 0
const GroupID = 14

var keyAPub []byte = nil
var keyBPub []byte = nil

func Send(side string, data []byte) bool {
	if(side == "Alice") {
		keyBPub = data
	} else if(side == "Bob") {
		keyAPub = data
	} else {
		log.Println("[ERROR] Invalid Side:", side)
		return false
	}
	log.Println("[DATA] Send to:", side, smart.Base64Encode(string(data)))
	return true
}

func Recv(side string) []byte {
	var out []byte = nil
	if(side == "Alice") {
		out = keyAPub
	} else if(side == "Bob") {
		out = keyBPub
	} else {
		log.Println("[ERROR] Invalid Side:", side)
		return nil
	}
	log.Println("[DATA] Recv from:", side)
	if(out == nil) {
		log.Println("[ERROR] Recv Side", side, "Key is NULL")
	}
	return out
}

func AliceSideStep1() (bool, *dhkx.DHGroup, *dhkx.DHKey) {

	// Get a group. Use the default one would be enough.
	g, errGrp := dhkx.GetGroup(GroupID)
	if(errGrp != nil) {
		log.Println("[ERROR] Get Group", "Alice", errGrp)
		return false, nil, nil
	}

	// Generate a private key from the group.
	// Use the default random number generator.
	priv, errGen := g.GeneratePrivateKey(nil)
	if(errGen != nil) {
		log.Println("[ERROR] Generate Private Key", "Alice", errGen)
		return false, nil, nil
	}
	if(priv == nil) {
		log.Println("[ERROR] Private Key is NULL", "Alice")
		return false, nil, nil
	}
	if(!priv.IsPrivateKey()) {
		log.Println("[ERROR] Private key is wrong", "Alice")
		return false, nil, nil
	}

	// Get the public key from the private key.
	pub := priv.Bytes()

	// Send the public key to Bob.
	ok := Send("Bob", pub)
	if(!ok) {
		log.Println("[ERROR] Wrong Answer to", "Alice", "from:", "Bob")
		return false, nil, nil
	}

	return true, g, priv

}

func AliceSideStep2(g *dhkx.DHGroup, priv *dhkx.DHKey) []byte {

	// Receive a slice of bytes from Bob, which contains Bob's public key
	b := Recv("Bob")

	// Recover Bob's public key
	bobPubKey := dhkx.NewPublicKey(b)
	log.Println("[DEBUG] Alice Side: Bob's Pub Key is:", smart.Base64Encode(string(bobPubKey.Bytes())))

	// Compute the key
	k, err := g.ComputeKey(bobPubKey, priv)
	if(err != nil) {
		log.Println("[ERROR] Compute Key", "Alice", err)
		return nil
	}

	// Get the key in the form of []byte
	key := k.Bytes()

	return key

}

func BobSideStep1() (bool, *dhkx.DHGroup, *dhkx.DHKey) {

	// Get a group. Use the default one would be enough.
	g, errGrp := dhkx.GetGroup(GroupID)
	if(errGrp != nil) {
		log.Println("[ERROR] Get Group", "Bob", errGrp)
		return false, nil, nil
	}

	// Generate a private key from the group.
	// Use the default random number generator.
	priv, errGen := g.GeneratePrivateKey(nil)
	if(errGen != nil) {
		log.Println("[ERROR] Generate Private Key", "Bob", errGen)
		return false, nil, nil
	}
	if(priv == nil) {
		log.Println("[ERROR] Private Key is NULL", "Bob")
		return false, nil, nil
	}
	if(!priv.IsPrivateKey()) {
		log.Println("[ERROR] Private key is wrong", "Bob")
		return false, nil, nil
	}

	// Get the public key from the private key.
	pub := priv.Bytes()

	// Send the public key to Alice.
	ok := Send("Alice", pub)
	if(!ok) {
		log.Println("[ERROR] Wrong Answer to", "Bob", "from:", "Alice")
		return false, nil, nil
	}

	return true, g, priv

}

func BobSideStep2(g *dhkx.DHGroup, priv *dhkx.DHKey) []byte {

	// Receive a slice of bytes from Alice, which contains Alice's public key
	b := Recv("Alice")

	// Recover Alice's public key
	alicePubKey := dhkx.NewPublicKey(b)
	log.Println("[DEBUG] Bob Side: Alice's Pub Key is:", smart.Base64Encode(string(alicePubKey.Bytes())))

	// Compute the key
	k, err := g.ComputeKey(alicePubKey, priv)
	if(err != nil) {
		log.Println("[ERROR] Compute Key", "Bob", err)
		return nil
	}

	// Get the key in the form of []byte
	key := k.Bytes()

	return key

}

/*
func BobSide() []byte {

	// Get a group. Use the default one would be enough.
	g, _ := dhkx.GetGroup(GroupID)

	// Generate a private key from the group.
	// Use the default random number generator.
	priv, errG := g.GeneratePrivateKey(nil)
	if(errG != nil) {
		log.Println("[ERROR] Generate Private Key", "Bob", errG)
		return nil
	}
	if(priv == nil) {
		log.Println("[ERROR] Private Key is NULL", "Bob")
		return nil
	}
	if(!priv.IsPrivateKey()) {
		log.Println("[ERROR] Private key is wrong", "Bob")
		return nil
	}

	// Get the public key from the private key.
	pub := priv.Bytes()

	// Send the public key to Alice.
	ok := Send("Alice", pub)
	if(!ok) {
		log.Println("[ERROR] Wrong Answer to", "Bob", "from:", "Alice")
		return nil
	}

	// Receive a slice of bytes from Alice, which contains Alice's public key
	a := Recv("Alice")

	// Recover Alice's public key
	alicePubKey := dhkx.NewPublicKey(a)
	log.Println("[DEBUG] Bob Side: Alice's Pub Key is:", smart.Base64Encode(string(alicePubKey.Bytes())))

	// Compute the key
	k, err := g.ComputeKey(alicePubKey, priv)
	if(err != nil) {
		log.Println("[ERROR] Compute Key:", err)
		return nil
	}

	// Get the key in the form of []byte
	key := k.Bytes()

	return key

}
*/


func main() {

	defer smart.PanicHandler()

	LogToConsoleWithColors()

	okA, gA, privA := AliceSideStep1()
	if(okA != true) {
		log.Println("[ERROR] Alice Step 1 Failed")
		return
	}
	okB, gB, privB := BobSideStep1()
	if(okB != true) {
		log.Println("[ERROR] Bob Step 1 Failed")
		return
	}
	keyA := AliceSideStep2(gA, privA)
	keyB := BobSideStep2(gB, privB)

	k64A := smart.Base64Encode(string(keyA))
	k64B := smart.Base64Encode(string(keyB))

	// To this point, the variables `key` on both Alice and Bob side are same. It could be used as the secret key for the later communication.

	log.Println("[INFO] Allice's computed Key:", k64A)
	log.Println("[INFO] Bob's computed Key:", k64B)

	if(k64A != k64B) {
		log.Println("[ERROR] Keys are different")
		return
	}
	log.Println("[OK] Keys are similar")

}

// #end
