package BLC

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

//网络服务文件管理

//3000作为引导节点（主节点）的地址
var knownNodes=[]string{"localhost:3000"}

//节点地址
var nodeAddress string

//启动服务
func startServer(nodeID string)  {
	fmt.Printf("启动服务[%s]...\n",nodeID)
	//节点地址赋值
	nodeAddress=fmt.Sprintf("localhost:%s",nodeID)
	//1.监听节点
	listen,err:=net.Listen(PROTOCOL,nodeAddress)
	if nil!=err{
		log.Panicf("listen address of %s failed!%v\n",nodeAddress,err)
	}
	defer listen.Close()
	//获取区块链blockchain对象
	bc:=BlockChainObject(nodeID)
	//两个节点，主节点负责保存数据，签到系欸但负责发送请求，同步数据
	if nodeAddress!=knownNodes[0]{
		//不是主节点，发送请求，同步数据
		//...
		//sendMessage(knownNodes[0],nodeAddress)
		sendVersion(knownNodes[0],bc)
	}

	for {
		//2.生成链接，接收请求
		conn,err:=listen.Accept()
		if nil!=err{
			log.Panicf("accept connect failed! %v\n",err)
		}
		//处理请求
		//单独启动一个goroutine进行请求处理
		go handleConnection(conn,bc)
	}
}
//请求处理函数
func handleConnection(conn net.Conn,bc *BlockChain)  {
	request,err:=ioutil.ReadAll(conn)
	if nil!=err{
		log.Panicf("Receive a Request failed! %v\n",err)
	}
	cmd:=bytesToCommand(request[:12])
	fmt.Printf("Receive a Command: %s\n",cmd)
	switch cmd {
	case CMD_VERSION :
		handleVersion(request,bc)
	case CMD_GETDATA:
		handleGetData(request,bc)
	case CMD_GETBLOCKS:
		handleGetBlocks(request,bc)
	case CMD_INV:
		handleInv(request,bc)
	case CMD_BLOCK:
		handleBlock(request,bc)
	default:
		fmt.Println("Unknow command")
	}
}



