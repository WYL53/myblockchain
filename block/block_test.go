package block

import (
	"testing"
)

func TestBlock(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("hello block chain.")
	bc.AddBlock("new block")

	for _,b := range bc.blocks {
		t.Logf("pre hash:%x\n",b.PrevBlockHash)
		t.Logf("data: %s\n",b.Data)
		t.Logf("hash:%x\n\n",b.Hash)
	}
}

