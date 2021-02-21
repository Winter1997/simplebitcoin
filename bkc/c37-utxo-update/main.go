package main

import "bitcoin/bkc/c37-utxo-update/BLC"

//启动
func main(){
	cli:=BLC.CLI{}
	cli.Run()
}