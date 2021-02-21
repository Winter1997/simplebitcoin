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
	//初始化区块连
	fmt.Println("\tcreateblockchain -address address ----创建区块链")
	//添加区块
	fmt.Println("\taddblock -data DATA ----添加区块")
	//打印完整的区块信息
	fmt.Println("\tprintchain----输出区块链信息")
	//通过命令转账
	fmt.Println("\t-from From -to -amount AMOUNT -- 发起转账")
	fmt.Println("\t转账参数说明")
	fmt.Println("\t\t-from FROM -- 转账源地址")
	fmt.Println("\t\t-to TO -- 转账目标地址")
	fmt.Println("\t\t-AMOUNT amount -- 转账金额")
}
//参数数量的检测函数
func IsValidArgs()  {
	if len(os.Args)<2{
		PrintUsage()
		//直接退出
		os.Exit(1)
	}
}
//发起交易
func (cli *CLI)send()  {

}


//初始化区块链
func (cli *CLI)createBlockchain(address string) {
	CreateBlockChainWithGenesisBlock(address)
}
//添加区块
func (cli *CLI)addBlock(txs []*Transaction)  {
	//判断数据库是否存在
	if !dbExist(){
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	blockchain:=BlockChainObject()
	//获取到blockchain的对象实例
	blockchain.AddBlock(txs)
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
	//发起交易
	sendCmd:=flag.NewFlagSet("send",flag.ExitOnError)
	//添加区块
	flagAddBlockArg:=addBlockCmd.String("data","sent 100 btc to player","添加区块数据")
	//创建区块链时指定的矿工地址（接收奖励）
	flagCreateBlockchainArg:=createBLCWithGenesisBlockCmd.String("address","troytan","指定接收系统奖励的矿工地址")
	//判断命令
	flagSendFromArg:=sendCmd.String("from","","转账源地址")
	flagSendToArg:=sendCmd.String("to","","转账目标地址")
	flagSendAmountArg:=sendCmd.String("amount","","转账金额")
	switch os.Args[1] {
	case "send":
		if err:=sendCmd.Parse(os.Args[2:]);nil!=err{
			log.Printf("parse sendCmd failed! %v\n",err)
		}
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
	//添加转账
	if sendCmd.Parsed(){
		if *flagSendFromArg==""{
			fmt.Printf("源地址不能为空...")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendToArg==""{
			fmt.Printf("目标地址不能为空...")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendAmountArg==""{
			fmt.Printf("转账金额不能为空...")
			PrintUsage()
			os.Exit(1)
		}
		fmt.Printf("\tFROM:[%s]\n",*flagSendFromArg)
		fmt.Printf("\tTO:[%s]\n",*flagSendToArg)
		fmt.Printf("\tAMOUNT:[%s]\n",*flagSendAmountArg)
	}
	//添加区块命令
	if addBlockCmd.Parsed(){
		if *flagAddBlockArg==""{
			PrintUsage()
			os.Exit(1)
		}
		cli.addBlock([]*Transaction{})
	}
	//输出区块链信息
	if printChainCmd.Parsed(){
		cli.printchain()
	}
	//创建区块链命令
	if createBLCWithGenesisBlockCmd.Parsed(){
		if *flagCreateBlockchainArg==""{
			PrintUsage()
			os.Exit(1)
		}
		cli.createBlockchain(*flagCreateBlockchainArg)
	}
}