package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

//对blockchain的命令进行管理

//client对象
type  CLI struct{
	BC *BlockChain  //区块链对象
}

//用法展示
func PrintUsage()  {
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain----创建区块链")
	fmt.Println("\taddblock -data DATA ----添加区块")
	fmt.Println("\tprintchain----输出区块链信息")
}

//初始化区块链
func (cli *CLI)createBlockchain() {
	CreateBlockChainWithGenesisBlock()
}
//添加区块
func (cli *CLI)addBlock(data string)  {
	//判断数据库是否存在
	if !dbExist(){
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	blockchain:=BlockChainObject()
	//获取到blockchain的对象实例
	blockchain.AddBlock([]byte(data))
	//cli.BC.AddBlock([]byte(data))
}
//打印完整区块链信息
func (cli *CLI)printchain()  {
	//判断数据库是否存在
	if !dbExist(){
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	blockchain:=BlockChainObject()
	blockchain.PrintChain()
}

//参数数量的检测函数
func IsValidArgs()  {
	if len(os.Args)<2{
		PrintUsage()
		//直接退出
		os.Exit(1)
	}
}

//命令行运行函数
func (cli *CLI)Run()  {
	//检测参数数量
	IsValidArgs()
	//新建相关命令
	//添加区块
	addBlockCmd:=flag.NewFlagSet("addblock",flag.ExitOnError)
	//输出区块链完整的信息
	printChainCmd:=flag.NewFlagSet("printchain",flag.ExitOnError)
	//创建区块链
	createBLCWithGenesisBlockCmd:=flag.NewFlagSet("createblockchain",flag.ExitOnError)
	//数据参数
	flagAddBlockArg:=addBlockCmd.String("data","sent 100 btc to player","添加区块数据")
	//判断命令
	switch os.Args[1] {
	case "addblock":
		if err:=addBlockCmd.Parse(os.Args[2:]);nil!=err{
			log.Panicf("parse addClockCmd failed! %v\n",err)
		}
	case "printchain":
		if err:=printChainCmd.Parse(os.Args[2:]);nil!=err{
			log.Panicf("parse printchainCmd failed! %v\n",err)
		}
	case "createblockchain":
		if err:=createBLCWithGenesisBlockCmd.Parse(os.Args[2:]);nil!=err{
			log.Panicf("parse createBLCWithGenesisBlockCmd failed! %v\n",err)
		}
	default:
		PrintUsage()
		os.Exit(1)
	}
	//添加区块命令
	if addBlockCmd.Parsed(){
		if *flagAddBlockArg==""{
			PrintUsage()
			os.Exit(1)
		}
		cli.addBlock(*flagAddBlockArg)
	}
	//输出区块链信息
	if printChainCmd.Parsed(){
		cli.printchain()
	}
	//创建区块链命令
	if createBLCWithGenesisBlockCmd.Parsed(){
		cli.createBlockchain()
	}
}