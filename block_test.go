package main

import (
	"testing"
	"strconv"
)

func TestBlock(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("hello block chain.")
	bc.AddBlock("new block")

	for _,b := range bc.blocks {
		t.Logf("pre hash:%x\n",b.PrevBlockHash)
		t.Logf("data: %s\n",b.Data)
		t.Logf("hash:%x\n",b.Hash)
		t.Logf("nonce:%x\n",b.Nonce)
		pow := NewProofOfWork(b)
		t.Logf("PoW : %s \n\n",strconv.FormatBool(pow.Validate()))
	}
}

