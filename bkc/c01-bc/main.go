package main

import (
	"bitcoin/bkc/c01-bc/BLC"
	"fmt"
)

//启动
func main(){
	block:=BLC.NewBlock(1,nil,[]byte("the first block testing"))
	fmt.Printf("the first block: %v\n",block)
}