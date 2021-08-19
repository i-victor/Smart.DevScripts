
package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"time"
	"strconv"
)

const (
	difficulty = 8
)

type Block struct {
	Index     uint64
	Timestamp int64
	Hash      []byte
	Data      string
	PrevHash  []byte
}

// bytes manipulation based on https://github.com/Jeiwan/blockchain_go

func genhash(data string, prev []byte) []byte {
	head := bytes.Join([][]byte{prev, []byte(data)}, []byte{})
	hdr := sha256.Sum256(head)
	fmt.Printf("Header hash: %x\n", hdr)
	return hdr[:]
}

func mine(hash []byte) []byte {
	target := big.NewInt(1)
	target = target.Lsh(target, uint(256 - difficulty))

	fmt.Printf("target: %x\n", target)

	var nonce int64

	for nonce = 0; nonce < math.MaxInt64; nonce++ {
		testNum := big.NewInt(0)
		testNum.Add(testNum.SetBytes(hash), big.NewInt(nonce))
		testHash := sha256.Sum256(testNum.Bytes())

		fmt.Printf("\rproof: %x (nonce: %d)", testHash, nonce)

		if target.Cmp(testNum.SetBytes(testHash[:])) > 0 {
			fmt.Println("\nFound!")
			return testHash[:]
		}
	}

	return []byte{}
}

func NewBlock(id uint64, data string, prev []byte) *Block {
	return &Block{
		id,
		time.Now().Unix(),
		mine(genhash(data, prev)),
		data,
		prev,
	}
}

func main() {

	start := time.Now()

	prev := []byte{}

	var i uint64 = 0
	var max uint64 = 1000
	var d string = ""
	for i = 0; i < max; i++ {
		if(i == 0) {
			d = "Genesis Block #0"
		} else {
			d = "Block #" + strconv.FormatUint(i, 10)
		}
		b := NewBlock(i, d, prev)
		fmt.Printf("Id: %d\nHash; %x\nData: %s\nPrevious: %x\n\n",
			b.Index,
			b.Hash,
			b.Data,
			b.PrevHash,
		)
		prev = b.Hash
	}

	duration := time.Since(start)

	fmt.Println("Total Time:", duration)

}

// #END
