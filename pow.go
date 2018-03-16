package main

import (
	"math/big"
	"bytes"
	"strconv"
	"fmt"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/common/math"
)

const targetBits = 24
const maxNonce  = math.MaxInt64

type ProofOfWork struct {
	block *Block
	target *big.Int //这是目标，需要找到比它小的hash值
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target,uint(256-targetBits)) // 1 << (256-targetBits)
	pow := &ProofOfWork{b,target}
	return pow
}

func (pow *ProofOfWork)prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.Data,
		IntToHex(pow.block.Timestamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	},[]byte{})
	return data
}

func IntToHex(x int64) []byte{
	return []byte(strconv.FormatInt(x,16))
}

func (pow *ProofOfWork) Run() (int,[]byte)  {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("mining the block containing \"%s\"\n",pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("nonce : %x\nhash:%x\n",nonce,hash)
			break
		}else {
			nonce++
		}
	}

	return nonce,hash[:]
}

func (pow *ProofOfWork)Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
