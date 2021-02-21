package BLC

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

//请求发送文件

//发送请求
func sendMessage(to string,msg []byte)  {
	fmt.Println("向服务器发送请求...")
	//1.连接上服务器
	conn,err:=net.Dial(PROTOCOL,to)
	if nil!=err{
		log.Panicf("connect to server [%s] failed!%v\n",err)
	}
	defer conn.Close()
	//要发送的数据
	_,err=io.Copy(conn,bytes.NewReader(msg))
	if nil!=err{
		log.Panicf("add the data to conn failed!%v\n",err)
	}
}

//区块链版本验证
func sendVersion(toAddress string,bc *BlockChain)  {
	//1.获取当前节点的区块高度
	height:=bc.GetHeight()
	//2.组装生成version
	versionData:=Version{Height: int(height),AddrFrom: nodeAddress}
	//3.数组序列化
	data:=gobEncode(versionData)
	//4.将命令与版本组装成完整的请求
	request:=append(commandToBytes(CMD_VERSION),data...)
	//5.发送请求
	sendMessage(toAddress,request)

}

//从指定区块同步数据
func sendGetBlocks(toAddress string)  {
	//1.生成数据
	data:=gobEncode(GetBlocks{AddrFrom: nodeAddress})
	//2.组装请求
	request:=append(commandToBytes(CMD_GETBLOCKS),data...)
	//3.发送请求
	sendMessage(toAddress,request)

}

//发送获取指定区块请求
func sendGetData(toAddress string,hash []byte)  {
	//1.生成数据
	data:=gobEncode(GetData{AddrFrom: nodeAddress,ID: hash})
	//2.组装请求
	request:=append(commandToBytes(CMD_GETDATA),data...)
	//3.发送请求
	sendMessage(toAddress,request)
}

//向其他节点展示
func sendInv(toAddress string,hashes [][]byte)  {
	//1.生成数据
	data:=gobEncode(Inv{AddrFrom: nodeAddress,Hashes: hashes})
	//2.组装请求
	request:=append(commandToBytes(CMD_INV),data...)
	//3.发送请求
	sendMessage(toAddress,request)
}

//发送区块信息
func sendBlock(toAddress string,block []byte)  {
	//1.生成数据
	data:=gobEncode(BlockData{AddrFrom: nodeAddress,Block: block})
	//2.组装请求
	request:=append(commandToBytes(CMD_BLOCK),data...)
	//3.发送请求
	sendMessage(toAddress,request)
}
