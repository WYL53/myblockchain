package main

import (
	"github.com/gpmgo/gopm/modules/log"
	"encoding/hex"
)

//交易输出
type TXOutpus struct {
	Value int
	ScriptPubKey string
}


//发送币
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInpus
	var outputs []TXOutpus

	acc,validOutputs := bc.FindSpendableOutputs(from,amount)
	if acc < amount{
		log.Warn("ERROR:Not enough funds")
	}

	//构建输入列表
	for txid,outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		for _,out := range outs{
			input := TXInput{txID,out,from}
			inpus - append(inputs,input)
		}
	}

	//构建输出列表
	outputs = append(outputs,TXOutpus{amount,to})
	if acc > amount{
		outputs = append(outputs,TXOutpus{acc-amount,from}) //找零
	}

	tx := Transaction{nil,inputs,outputs}
	tx.SetID()
	return &tx
}

type TXInput struct {
	Txid []byte //以前交易的ID
	Vout int  //该输出的索引
	ScriptSig string
}
