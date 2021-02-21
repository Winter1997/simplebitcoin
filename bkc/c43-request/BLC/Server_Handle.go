package BLC

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

//请求处理文件管理


//version
func handleVersion(request []byte,bc *BlockChain)  {
	fmt.Println("the request of version handle...")
	var buff bytes.Buffer
	var data Version
	//1.解析请求
	dataBytes:=request[12:]
	//2.生成version结构
	buff.Write(dataBytes)
	decoder:=gob.NewDecoder(&buff)
	if err:=decoder.Decode(&data);nil!=err{
		log.Panicf("decode the version struct failed! %v\n",err)
	}
	//3.获取请求方的区块高度
	versionHeight:=data.Height
	//4.获取自身节点的区块高度
	height:=bc.GetHeight()
	fmt.Printf("height : %v ,versionHeight : %v\n",height,versionHeight)
	//如果当前节点的区块高度大于versionHeight
	//将当前节点版本信息发送给请求节点
	if height>int64(versionHeight){
		sendVersion(data.AddrFrom,bc)

	}else if height<int64(versionHeight){
		//如果当前节点区块高度小于versionHeight
		//向发送方发起同步数据的请求
		sendGetBlocks(data.AddrFrom)
	}
}

//GetBlocks
//数据同步请求处理
func handleGetBlocks(request []byte,bc *BlockChain)  {
	fmt.Println("the request of get blocks handle...")
	var buff bytes.Buffer
	var data GetBlocks
	//1.解析请求
	dataBytes:=request[12:]
	//2.生成getblocks结构
	buff.Write(dataBytes)
	decoder:=gob.NewDecoder(&buff)
	if err:=decoder.Decode(&data);nil!=err{
		log.Panicf("decode getblocks struct failed! %v\n",err)
	}
	//3.获取区块链所有的区块哈希
	hashes:=bc.GetBlockHases()
	sendInv(data.AddrFrom,hashes)
}

//Inv
func handleInv(request []byte,bc *BlockChain)  {
	fmt.Println("the request of inv handle...")
	var buff bytes.Buffer
	var data Inv
	//1.解析请求
	dataBytes:=request[12:]
	//2.生成Inv结构
	buff.Write(dataBytes)
	decoder:=gob.NewDecoder(&buff)
	if err:=decoder.Decode(&data);nil!=err{
		log.Panicf("decode inv struct failed! %v\n",err)
	}
	sendGetData(data.AddrFrom,data.Hashes[0])
	/*for _,hash:=range data.Hashes{
		sendGetData(data.AddrFrom,hash)
	}*/
}

//GetData
//处理获取指定区块的请求
func handleGetData(request []byte,bc*BlockChain)  {
	fmt.Println("the request of get block handle...")
	var buff bytes.Buffer
	var data GetData
	//1.解析请求
	dataBytes:=request[12:]
	//2.生成getData结构
	buff.Write(dataBytes)
	decoder:=gob.NewDecoder(&buff)
	if err:=decoder.Decode(&data);nil!=err{
		log.Panicf("decode getData struct failed! %v\n",err)
	}
	//3.通过传过来的区块哈希，获取本地节点的区块
	blockBytes:=bc.GetBlock(data.ID)
	sendBlock(data.AddrFrom,blockBytes)
}

//Block
//接收到新区块的时候进行处理
func handleBlock(request []byte,bc *BlockChain)  {
	fmt.Println("the request of handle block handle...")
	var buff bytes.Buffer
	var data BlockData
	//1.解析请求
	dataBytes:=request[12:]
	//2.生成blockdata结构
	buff.Write(dataBytes)
	decoder:=gob.NewDecoder(&buff)
	if err:=decoder.Decode(&data);nil!=err{
		log.Panicf("decode blockdata struct failed! %v\n",err)
	}
	//3.将接收到的区块添加到区块链中
	blockBytes:=data.Block
	block:=DeserializeBlock(blockBytes)
	bc.AddBlock(block)
	//4.更新utxo table
	utxoSet:=UTXOSet{bc}
	utxoSet.update()
}