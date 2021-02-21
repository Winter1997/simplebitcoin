package main

import "bitcoin/bkc/c35-utxo-table-op/BLC"

//启动
func main(){
	cli:=BLC.CLI{}
	cli.Run()
}