package BLC

//公共模块

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

//参数数量的检测函数
func IsValidArgs()  {
	if len(os.Args)<2{
		PrintUsage()
		//直接退出
		os.Exit(1)
	}
}
//实现int64转[]byte
func IntToHex(data int64) []byte {
	buffer:=new(bytes.Buffer)
	err:=binary.Write(buffer,binary.BigEndian,data)
	if err!=nil{
		log.Panicf("int transact to []byte failed! %v\n",err)
	}
	return buffer.Bytes()
}

//标准JSON格式转切片
//windows下需要添加引号
//bc.exe send -from "[\"troytan\"]" -to "[\"Alice\"]" -amount "[\"5\"]"
//bc.exe send -from "[\"troytan\",\"Alice\"]" -to "[\"Alice\",\"troytan\"]" -amount "[\"5\",\"2\"]"
//troytan ->Alice 5 -->Alice 5 troytan 5
//Alice ->troytan 2 -->Alice 3 troytan 7
func JSONToSlice(jsonString string) []string {
	var strSlice []string
	//json
	if err:=json.Unmarshal([]byte(jsonString),&strSlice);nil!=err{
		log.Panicf("json to []string fqailed! %v\n",err)
	}
	return strSlice
}
//string 转hash160
func StringToHash160(address string) []byte {
	pubKeyHash:=Base58Decode([]byte(address))
	hash160:=pubKeyHash[:len(pubKeyHash)-addressCheckSumLen]
	return hash160
}
//获取节点ID
func GetEnvNodeId() string {
	nodeID:=os.Getenv("NODE_ID")
	if nodeID==""{
		fmt.Println("NODE_ID is not set...")
		os.Exit(1)
	}
	return nodeID
}

//gob编码
func gobEncode(data interface{}) []byte {
	var result bytes.Buffer
	enc:=gob.NewEncoder(&result)
	err:=enc.Encode(data)
	if nil!=err{
		log.Printf("encode the data failed! %v\n",err)
	}
	return result.Bytes()
}

//命令转换为请求（[]byte）
func commandToBytes(command string) []byte {
	var bytes [CMMAND_LENGTH]byte
	for i,c:=range command{
		bytes[i]=byte(c)
	}
	return bytes[:]
}

