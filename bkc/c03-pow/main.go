package main

import (
	"bitcoin/bkc/c03-pow/BLC"
	"fmt"
)

//启动
func main(){
	bc:=BLC.CreateBlockChainWithGenesisBlock()
	fmt.Printf("blockchain:%v\n",bc.Blocks[0])
	//上链
	bc.AddBlock(bc.Blocks[len(bc.Blocks)-1].Height+1,bc.Blocks[len(bc.Blocks)-1].Hash,[]byte("alice send 10 btc to bob"))
	bc.AddBlock(bc.Blocks[len(bc.Blocks)-1].Height+1,bc.Blocks[len(bc.Blocks)-1].Hash,[]byte("bob send 5 to troytan"))
	for _,block:=range bc.Blocks{
		fmt.Printf("prevBlockHash : %x, currentHash : %x\n",block.PreBlockHash,block.Hash)
	}
}