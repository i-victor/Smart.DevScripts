package main

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce int64 = math.MaxInt64
)

const targetBits = 16

// ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(512 - targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Run performs a proof-of-work
func (pow *ProofOfWork) Run() (int64, []byte) {
	var hashInt big.Int
	var hash [64]byte
	var nonce int64 = 0

//	fmt.Printf("Mining a new block")
	fmt.Println("******** Mining a new block ********")
	for nonce < maxNonce {
		data := pow.prepareData(nonce)

		hash = sha512.Sum512(data)
		if math.Remainder(float64(nonce), 100000) == 0 {
			fmt.Printf("\r%x", hash)
		}
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate validates block's PoW
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha512.Sum512(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
