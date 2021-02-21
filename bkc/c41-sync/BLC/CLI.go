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
	//初始化区块链
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
	//查询余额
	fmt.Printf("\tgetbalance -address FROM -- 查询指定地址的余额\n")
	fmt.Println("\t查询余额参数说明")
	fmt.Printf("\t\t-address -- 查询余额的地址\n")
	//钱包管理
	fmt.Println("\tcreatewallet --创建钱包")
	fmt.Println("\taccounts --获取地址列表")
	fmt.Printf("\tutxo -method METHOD -- 测试UTXO Table功能中指定的方法\n")
	fmt.Printf("\t\tMETHOD -- 方法名\n")
	fmt.Printf("\t\t\treset -- 重置UTXOtable\n")
	fmt.Printf("\t\t\tbalance - 查找所有UTXO\n")
	fmt.Printf("\t\tport -- 访问的节点号\n")
	fmt.Printf("\tstart -- 启动节点服务\n")
}

//添加区块
func (cli *CLI)addBlock(txs []*Transaction)  {
	//判断数据库是否存在
	//if !dbExist(){
	//	fmt.Println("数据库不存在...")
	//	os.Exit(1)
	//}
	//blockchain:=BlockChainObject()
	//获取到blockchain的对象实例
	//blockchain.AddBlock(txs)
	//cli.BC.AddBlock([]byte(data))
}

//命令行运行函数
func (cli *CLI)Run()  {
	nodeId:=GetEnvNodeId()
	//检测参数数量
	IsValidArgs()
	//新建相关命令
	//添加区块
	addBlockCmd:=flag.NewFlagSet("addblock",flag.ExitOnError)
	//创建钱包
	createWalletCmd:=flag.NewFlagSet("createwallet",flag.ExitOnError)
	//获取地址列表
	getAccountsCmd:=flag.NewFlagSet("getaccounts",flag.ExitOnError)
	UTXOTestCmd:=flag.NewFlagSet("utxo",flag.ExitOnError)
	//输出区块链完整的信息
	printChainCmd:=flag.NewFlagSet("printchain",flag.ExitOnError)
	//创建区块链
	createBLCWithGenesisBlockCmd:=flag.NewFlagSet("createblockchain",flag.ExitOnError)
	//节点号设置命令
	setNodeIdCmd:=flag.NewFlagSet("set_id",flag.ExitOnError)
	//节点服务启动命令
	startNodeCmd:=flag.NewFlagSet("start",flag.ExitOnError)
	//数据参数
	//发起交易
	sendCmd:=flag.NewFlagSet("send",flag.ExitOnError)
	//查询余额的命令
	getBalanceCmd:=flag.NewFlagSet("getBalance",flag.ExitOnError)
	//添加区块
	flagAddBlockArg:=addBlockCmd.String("data","sent 100 btc to player","添加区块数据")
	//创建区块链时指定的矿工地址（接收奖励）
	flagCreateBlockchainArg:=createBLCWithGenesisBlockCmd.String("address","troytan","指定接收系统奖励的矿工地址")
	//UTXO测试命令行参数
	flagUTXOArg:=UTXOTestCmd.String("method","","UTXO Table相关操作\n")
	//端口号参数
	flagPortArg:=setNodeIdCmd.String("port","","设置节点 ID")
	//判断命令
	flagSendFromArg:=sendCmd.String("from","","转账源地址")
	flagSendToArg:=sendCmd.String("to","","转账目标地址")
	flagSendAmountArg:=sendCmd.String("amount","","转账金额")
	flagGetBalanceArg:=getBalanceCmd.String("address","","要查询的地址")
	switch os.Args[1] {
	case "start":
		if err:=startNodeCmd.Parse(os.Args[2:]);nil!=err{
			log.Panicf("parse cmd start node server failed! %v\n",err)
		}
	case "set_id":
		if err:=setNodeIdCmd.Parse(os.Args[2:]);nil!=err{
			log.Panicf("parse cmd set node id failed! %v\n",err)
		}
	case "utxo":
		if err:=UTXOTestCmd.Parse(os.Args[2:]);nil!=err{
			log.Panicf("parse cmd utxo table failed! %v\n",err)
		}
	case "accounts":
		if err:=getAccountsCmd.Parse(os.Args[2:]);nil!=err{
			log.Panicf("parse of cmd getaccounts failed! %v\n",err)
		}
	case "createwallet":
		if err:=createWalletCmd.Parse(os.Args[2:]);nil!=err{
			log.Panicf("parse cmd of create wallet failed! %v\n",err)
		}
	case "getbalance":
		if err:=getBalanceCmd.Parse(os.Args[2:]);nil!=err{
			log.Panicf("parse cmd get balance failed! %v\n",err)
		}
	case "send":
		if err:=sendCmd.Parse(os.Args[2:]);nil!=err{
			log.Panicf("parse sendCmd failed! %v\n",err)
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
	//节点启动服务
	if startNodeCmd.Parsed(){
		cli.startNode(nodeId)
	}
	//节点ID设置
	if setNodeIdCmd.Parsed(){
		if *flagPortArg==""{
			fmt.Printf("请输入端口号...")
			os.Exit(1)
		}
		cli.SetNodeId(*flagPortArg)
	}
	//utxo table操作
	if UTXOTestCmd.Parsed(){
		switch *flagUTXOArg {
		case "balance":
			cli.TestFindUTXOMap()
		case "reset":
			cli.TestResetUTXO(nodeId)
		default:

		}
	}
	//获取地址列表
	if getAccountsCmd.Parsed(){
		cli.GetAccounts(nodeId)
	}
	//创建钱包
	if createWalletCmd.Parsed(){
		cli.CreateWallets(nodeId)
	}
	//查询余额
	if getBalanceCmd.Parsed(){
		if *flagGetBalanceArg==""{
			fmt.Println("请输入查询地址...")
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceArg,nodeId)
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
		fmt.Printf("\tFROM:[%s]\n",JSONToSlice(*flagSendFromArg))
		fmt.Printf("\tTO:[%s]\n",JSONToSlice(*flagSendToArg))
		fmt.Printf("\tAMOUNT:[%s]\n",JSONToSlice(*flagSendAmountArg))
		cli.send(JSONToSlice(*flagSendFromArg),JSONToSlice(*flagSendToArg),JSONToSlice(*flagSendAmountArg),nodeId)
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
	fmt.Println("nodeID: ",nodeId)
	if printChainCmd.Parsed(){
		cli.printchain(nodeId)
	}
	//创建区块链命令
	if createBLCWithGenesisBlockCmd.Parsed(){
		if *flagCreateBlockchainArg==""{
			PrintUsage()
			os.Exit(1)
		}
		cli.createBlockchain(*flagCreateBlockchainArg,nodeId)
	}
}