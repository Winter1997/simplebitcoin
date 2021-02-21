package BLC

import (
	"bytes"
	"fmt"
	"io"
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
	listen,err:=net.Listen("tcp",nodeAddress)
	if nil!=err{
		log.Panicf("listen address of %s failed!%v\n",nodeAddress,err)
	}
	defer listen.Close()
	//两个节点，主节点负责保存数据，签到系欸但负责发送请求，同步数据
	if nodeAddress!=knownNodes[0]{
		//不是主节点，发送请求，同步数据
		//...
		sendMessage(knownNodes[0],nodeAddress)
	}

	for {
		//2.生成链接，接收请求
		conn,err:=listen.Accept()
		if nil!=err{
			log.Panicf("accept connect failed! %v\n",err)
		}
		request,err:=ioutil.ReadAll(conn)
		if nil!=err{
			log.Panicf("Receive Message failed! %v\n",err)
		}
		//3.处理请求
		fmt.Printf("Receive Message : %v\n",request)
	}
}
//发送请求
func sendMessage(to,from string)  {
	fmt.Println("向服务器发送请求...")
	//1.连接上服务器
	conn,err:=net.Dial("tcp",to)
	if nil!=err{
		log.Panicf("connect to server [%s] failed!%v\n",err)
	}


	//要发送的数据
	_,err=io.Copy(conn,bytes.NewReader([]byte(from)))
	if nil!=err{
		log.Panicf("add the data to conn failed!%v\n",err)
	}
}