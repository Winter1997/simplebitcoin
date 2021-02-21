package main

import "bitcoin/bkc/c42-request-gen/BLC"

//启动
func main(){
	cli:=BLC.CLI{}
	cli.Run()
}