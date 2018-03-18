package main

import (
	"flag"
	"os"
	"fmt"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}

func (cli *CLI) Run() {
	//cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)

	addBlockData := addBlockCmd.String("data","","Block data")
	switch os.Args[1] {
	case "addblock":
		addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		printChainCmd.Parse(os.Args[2:])
	default:
		//cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed(){
		if *addBlockData == ""{
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed(){
		cli.printChain()
	}
}

func (cli *CLI) addBlock(data string)  {
	cli.bc.AddBlock(data)
	fmt.Println("success")
}

func (cli *CLI) printChain()  {
	bci := cli.bc.Iterator()
	for  {
		block := bci.Next()
		fmt.Printf("prev hash:%x\n",block.PrevBlockHash)
		fmt.Printf("data:%s\n",block.Data)
		fmt.Printf("hash:%x\n",block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW:%s\n",strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PrevBlockHash) == 0{
			break
		}
	}
}


func (cli *CLI)send(from, to string, amount int)  {
	bc := NewBlockchain(from)
	defer bc.db.Close()

	tx := NewUTXOTransaction(from,to,amount,bc)
	bc.MineBlock([]*Transacton{tx})
	fmt.Println("success.")
}