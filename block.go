package main

import (
	"time"
	"bytes"
	"encoding/gob"
	"github.com/boltdb/bolt"
	"encoding/hex"
)

const dbFile  = "./dbFile"
const blocksBucket = "blocksBucket"

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
	tip []byte
	db *bolt.DB
}

func (bc *Blockchain)AddBlock(data string)  {
	var lastHash []byte
	 bc.db.View(func(tx *bolt.Tx)error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("1"))
		return nil
	})

	newBlock := NewBlock(data,lastHash)

	bc.db.Update(func(tx *bolt.Tx)error {
		b := tx.Bucket([]byte(blocksBucket))
		b.Put(newBlock.Hash,newBlock.Serialize())
		b.Put([]byte("1"),newBlock.Hash)
		bc.tip = newBlock.Hash
		return nil
	})
}

func (bc *Blockchain)FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

	//对所有未花费的交易进行迭代，进行累加，找到大于等于amount
	work:
	for _,tx := range unspentTXs{
		txID := hex.EncodeToString(tx.ID)
		for outIdx,out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentTXs[txID], outIdx)
				if accumulated >= amount {
					break work
				}
			}
		}
	}
	return accumulated,unspentOutputs
}

func (bc *Blockchain)MineBlock(transactions []*Transaction)  {

	newBlock := NewBlock(transactions,lastHash)
}

func (b  *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil{
		panic(err)
	}
	return result.Bytes()
}

func DeserializedBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil{
		panic(err)
	}
	return &block
}



func NewBlockchain() *Blockchain {
	var tip []byte
	db,err := bolt.Open(dbFile,0600,nil)
	if err != nil{
		panic("open db file failed.")

	}
	err = db.Update(func(tx *bolt.Tx)error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil{
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash,genesis.Serialize())
			b.Put([]byte("1"),genesis.Hash)
			tip = genesis.Hash
			return err
		}else {
			tip = b.Get([]byte("1"))
			return nil
		}

	})
	return &Blockchain{tip,db}
}

//创建创世区块
func NewGenesisBlock() *Block {
	return NewBlock("this is genesis block",[]byte{})
}

type BlockchainIterator struct {
	currentHash []byte
	db *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator  {
	bci := &BlockchainIterator{
		currentHash:bc.tip,
		db : bc.db,
	}
	return bci
}

func (bci *BlockchainIterator)Next() *Block {
	var block *Block

	bci.db.View(func(tx *bolt.Tx)error {
		b := tx.Bucket([]byte(blocksBucket))
		encodeBlock := b.Get(bci.currentHash)
		block = DeserializedBlock(encodeBlock)
		return nil
	})
	bci.currentHash = block.PrevBlockHash
	return block
}

