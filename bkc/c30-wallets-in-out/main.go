package main

import "bitcoin/bkc/c30-wallets-in-out/BLC"

//启动
func main(){
	cli:=BLC.CLI{}
	cli.Run()
}