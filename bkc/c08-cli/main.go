package main

import (
	"flag"
	"fmt"
)

//命令行回顾
//定义一个字符串变量
var species=flag.String("species","go","the usage of flag")

//定义一个int字符
var num =flag.Int("ins",1,"ins nums")
func main()  {
	//解析，在flags各类型参数生效之前，需要对参数进行解析
	flag.Parse()
	//打印参数
	fmt.Println("a string flag",*species)
	fmt.Println("ins num:",*num)
}