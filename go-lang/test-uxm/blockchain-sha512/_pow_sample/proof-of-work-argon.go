
// bytes manipulation based on https://github.com/Jeiwan/blockchain_go

package main

import (
	"os"
	"log"
	"bytes"
	"fmt"
	"math"
	"math/big"
	"time"
	"strconv"
//	"encoding/hex"
	"golang.org/x/crypto/argon2"

	color "github.com/unix-world/smartgo/colorstring"
	smart "github.com/unix-world/smartgo"
	b58   "github.com/unix-world/smartgo/base58"
	b92   "github.com/unix-world/smartgo/base92"
	uid   "github.com/unix-world/smartgo/uuid"
)

func LogToConsoleWithColors() {
	//--
	smart.ClearPrintTerminal()
	//--
//	smart.LogToStdErr("DEBUG")
	smart.LogToConsole("DEBUG", true) // colored or not
//	smart.LogToFile("DEBUG", "logs/", "json", true, true) // json | plain ; also on console ; colored or not
	//--
	log.Println("[DEBUG] Debugging")
	log.Println("[DATA] Data")
	log.Println("[NOTICE] Notice")
	log.Println("[WARNING] Warning")
	log.Println("[ERROR] Error")
	log.Println("[OK] OK")
	log.Println("A log message, with no type")
	//--
} //END FUNCTION


func fatalError(logMessages ...interface{}) {
	//--
	log.Fatal("[ERROR] ", logMessages) // standard logger
	os.Exit(1)
	//--
} //END FUNCTION

const (
	difficulty uint64 = 8 // 4 * zeroesPrefixNum
	argonsalt = "ProofOfWork:Argon2id"
)

func generateArgon2idHash(data []byte) []byte {
	return argon2.IDKey(data, []byte(argonsalt), 12, 16*1024, 1, 64) // 21 256 1 64
}

func hashToASCII(data []byte) string {
	return b92.Encode(data)
}

func ASCIItohash(data string) []byte {
	bytes, err := b92.Decode(data)
	if(err != nil) {
		return nil
	}
	return bytes
}

type Block struct {
	Id         string
	DateTime   string
	CheckSum   string
	Difficulty uint64
	Nonce      int64
	Hash       []byte
	PrevHash   []byte
	Data       string
}
type JsonBlock struct {
	Id         string `json:"id"`
	DateTime   string `json:"datetime"`
	CheckSum   string `json:"checksum"`
	Difficulty string `json:"difficulty"`
	Nonce      string `json:"nonce"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevhash"`
	Data       string `json:"data"`
}

func getDifficultyTarget() *big.Int {
	target := big.NewInt(1)
	target = target.Lsh(target, uint(512 - difficulty))
	return target
}

func verifyHash(target *big.Int, hdrHash []byte, nonce int64) (ok bool, bytes []byte) {
	testNum := big.NewInt(0)
	testNum.Add(testNum.SetBytes(hdrHash), big.NewInt(nonce))
	testHash := generateArgon2idHash(testNum.Bytes())
	if target.Cmp(testNum.SetBytes(testHash[:])) > 0 {
		return true, testHash[:]
	}
	return false, testHash[:]
}

func mine(data string, prevHash []byte) (theHash []byte, theNonce int64, theDifficulty uint64) {
	target := getDifficultyTarget()
	fmt.Println("\n########")
//	fmt.Printf("target: %x\n", target)
	log.Println("[DEBUG] Difficulty: " + smart.ConvertUInt64ToStr(difficulty))
	log.Println("[DEBUG] Target: " + hashToASCII(target.Bytes()))
	hdrHash := generateHeaderHash(data, prevHash)
//	fmt.Printf("Header-Hash: %x\n", hdrHash)
	fmt.Println("Header-Hash: " + hashToASCII(hdrHash))
	var nonce int64 = -1
	for nonce = 0; nonce < math.MaxInt64; nonce++ {
		testVerify, theBytes := verifyHash(target, hdrHash, nonce)
		fmt.Printf("\rProof-of-Work: %s (Nonce: %d)", hashToASCII(theBytes), nonce)
		if(testVerify == true) {
			fmt.Println("")
			log.Println("[NOTICE] Found ; NONCE=" + smart.ConvertInt64ToStr(nonce))
			return theBytes, nonce, difficulty
		}
	}
	log.Println("[WARNING] NOT Found")
	return []byte{}, -1, difficulty
}

func generateHeaderHash(data string, prev []byte) []byte {
	head := bytes.Join([][]byte{prev, []byte(data)}, []byte{})
	hdr := generateArgon2idHash(head)
	return hdr[:]
}

func generateBlockChecksum(id string, dTime string, data string, theNonce int64, theHash []byte, prevHash []byte) string {
	var strCk string = smart.RawUrlEncode(id) + "\n" + smart.RawUrlEncode(dTime) + "\n" + smart.Base64Encode(data) + "\n" + smart.ConvertInt64ToStr(theNonce) + "\n" + smart.Base64Encode(string(theHash)) + "\n" + smart.Base64Encode(string(prevHash))
//	fmt.Println(strCk)
	return b92.Encode([]byte(smart.Hex2Bin(smart.Sha512(strCk))))
}

func NewBlock(data string, prevHash []byte) *Block {
	dtObjUtc := smart.DateTimeStructUtc("")
	if(dtObjUtc.Status != "OK") {
		fatalError("Date Time ERROR: " + dtObjUtc.ErrMsg)
		return nil
	}
	theHash, theNonce, theDifficulty := mine(data, prevHash)
	theTime := smart.ConvertInt64ToStr(dtObjUtc.Time)
	var dt string = dtObjUtc.Years + "-" + dtObjUtc.Months + "-" + dtObjUtc.Days + " " + dtObjUtc.Hours + ":" + dtObjUtc.Minutes + ":" + dtObjUtc.Seconds + " " + dtObjUtc.TzOffset
	var id string = b58.Encode([]byte(theTime)) + "-" + uid.Uuid1013Str(13)
	return &Block{
		id,
		dt,
		generateBlockChecksum(id, dt, data, theNonce, theHash, prevHash),
		theDifficulty,
		theNonce,
		theHash,
		prevHash,
		data,
	}
}

func main() {

	LogToConsoleWithColors()

	start := time.Now()

	prev := []byte{}

	var i uint64 = 0
	var max uint64 = 3
	var d string = ""
	var b *Block
	for i = 0; i < max; i++ {
		if(i == 0) {
			d = "Genesis Block #0"
		} else {
			d = "Block #" + strconv.FormatUint(i, 10)
		}
		d = smart.JsonRawEncode(d) // ISSUE: json encode may vary ... and the block header will vary accordingly ; TODO: add data as base64#Line1 \n base64#line2 ...
		b = NewBlock(d, prev)
		fmt.Printf("Id: %s\nData: %s\nHash: %s\nNonce: %d\nPrevious: %s\n########\n",
			b.Id,
			b.Data,
			hashToASCII(b.Hash),
			b.Nonce,
			hashToASCII(b.PrevHash),
		)
		prev = b.Hash
	}

	duration := time.Since(start)
	fmt.Println("")
	log.Println("[INFO] ======== TOTAL Hashing Time:", duration, "========")
	fmt.Println("")

	target := getDifficultyTarget()

	testVerify, _ := verifyHash(target, generateHeaderHash(b.Data, b.PrevHash), b.Nonce)
	if(testVerify == true) {
		log.Println("[OK] BLOCK: Last Block Check is OK")
	} else {
		log.Println("[ERROR] BLOCK: Block Check is NOT OK !!!!!!!!")
	}
	log.Println("[DATA] Block Hex Hash:", hashToASCII(b.Hash))
	log.Println("[DATA] Block Prev Hex Hash:", hashToASCII(b.PrevHash))
	j := &JsonBlock {
		b.Id,
		b.DateTime,
		b.CheckSum,
		smart.ConvertUInt64ToStr(b.Difficulty),
		smart.ConvertInt64ToStr(b.Nonce),
		hashToASCII(b.Hash),
		hashToASCII(b.PrevHash),
		smart.Base64Encode(b.Data),
	}
	var jsonStr string = smart.JsonRawEncodePretty(j)
//	var jsonStr string = smart.JsonRawEncode(j)
	fmt.Println(color.YellowString("[DATA] LastBlock Data Json: " + jsonStr))

	D := smart.JsonDecode(jsonStr)
	if(D == nil) {
		log.Println("[ERROR] JSON: Decode Fail")
		return;
	}
	keys := []string{"id", "datetime", "checksum", "difficulty", "nonce", "hash", "prevhash", "data"}
	for y, k := range keys {
		if(D[k] == nil) {
			log.Println("[ERROR] JSON: Key is missing: #" + smart.ConvertIntToStr(y) + " @ " + k)
			return;
		}
	}
	B := &Block{
		D["id"].(string),
		D["datetime"].(string),
		D["checksum"].(string),
		smart.ParseIntegerStrAsUInt64(D["difficulty"].(string)),
		smart.ParseIntegerStrAsInt64(D["nonce"].(string)),
		ASCIItohash(D["hash"].(string)),
		ASCIItohash(D["prevhash"].(string)),
		smart.Base64Decode(D["data"].(string)),
	}

	log.Println("[DATA] Block Prev Hex Hash:", hashToASCII(B.PrevHash))
	log.Println("[DATA] Block Hex Hash:", hashToASCII(B.Hash))
	log.Println("[DATA] Block Data:", B.Data)
	testVerify2, _ := verifyHash(target, generateHeaderHash(B.Data, B.PrevHash), B.Nonce)
	if(testVerify2 == true) {
		log.Println("[OK] JSON: Last Block Check is OK")
	} else {
		log.Println("[ERROR] JSON: Last Block Check is NOT OK")
	}
	var ckSum string = generateBlockChecksum(B.Id, B.DateTime, B.Data, B.Nonce, B.Hash, B.PrevHash)
	if(ckSum == B.CheckSum) {
		log.Println("[OK] JSON: Checksum OK:", "`" + ckSum + "`")
	} else {
		log.Println("[ERROR] JSON: Checksum does not match:", "`" + ckSum + "`")
	}

}

// #END
