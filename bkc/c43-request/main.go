package main

import "bitcoin/bkc/c43-request/BLC"

//启动
func main(){
	cli:=BLC.CLI{}
	cli.Run()
}