package main

import (
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce 			int
}

//func (b *Block) SetHash() {
//	timeStamp := []byte(strconv.FormatInt(b.Timestamp, 10))
//	header := bytes.Join([][]byte{b.PrevBlockHash,b.Data,timeStamp},[]byte{})
//	hash := sha256.Sum256(header)
//	b.Hash = hash[:]
//}

func NewBlock(data string,prevBlockHash []byte) *Block  {
	block := &Block{
		Timestamp:time.Now().Unix(),
		Data:[]byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:nil,
		Nonce:0,
	}
	pow := NewProofOfWork(block)
	nonce,hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce

	return block
}

type Blockchain struct {
	blocks []*Block
}

func (bc *Blockchain)AddBlock(data string)  {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data,prevBlock.Hash)
	bc.blocks = append(bc.blocks,newBlock)
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

//创建创世区块
func NewGenesisBlock() *Block {
	return NewBlock("this is genesis block",[]byte{})
}
