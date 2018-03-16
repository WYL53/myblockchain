package main

import (
	"strconv"
	"fmt"
)

func main()  {
	bc := NewBlockchain()
	bc.AddBlock("hello block chain.")
	bc.AddBlock("new block")

	for _,b := range bc.blocks {
		fmt.Printf("pre hash:%x\n",b.PrevBlockHash)
		fmt.Printf("data: %s\n",b.Data)
		fmt.Printf("hash:%x\n",b.Hash)
		fmt.Printf("nonce:%x\n",b.Nonce)
		pow := NewProofOfWork(b)
		fmt.Printf("PoW : %s \n\n",strconv.FormatBool(pow.Validate()))
	}
}
