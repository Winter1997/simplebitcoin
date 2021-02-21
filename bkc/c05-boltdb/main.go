package main

import (
	"bitcoin/bkc/c05-boltdb/BLC"
	"fmt"
	"github.com/boltdb/bolt"
)

//启动
func main(){
	bc:=BLC.CreateBlockChainWithGenesisBlock()
	bc.AddBlock([]byte("a send 100 eth to b"))
	bc.AddBlock([]byte("b send 100 eth to a"))
	bc.DB.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte("blocks"))
		if nil!=b{
			hash:=b.Get([]byte("1"))
			fmt.Printf("value : %x\n",hash)
		}
		return nil
	})
}