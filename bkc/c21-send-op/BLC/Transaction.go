package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

//交易管理文件

//定义一个交易基本结构
type Transaction struct{
	//交易哈希（标识）
	TxHash   []byte
	//输入列表
	Vins     []*TxInput
	//输出列表
	Vouts    []*TxOutput
}

//实现coinbase交易
func NewCoinbaseTransaction(address string) *Transaction {
	//输入
	//coinbase
	//txHash:nil
	//vout:-1（为了对是否是coinbase交易进行判断）
	//ScriptSig:系统奖励
	txInput:=&TxInput{[]byte{},-1,"system reward"}
	//输出：
	//value：
	//address：
	txOuput:=&TxOutput{10,address}
	//输入输出组装交易
	txCoinbase:=&Transaction{nil,[]*TxInput{txInput},[]*TxOutput{txOuput}}
	//交易哈希生成
	txCoinbase.HashTransaction()
	return txCoinbase
}
//生成交易哈希（交易序列化）
func (tx *Transaction) HashTransaction()  {
	var result bytes.Buffer
	//设置编码对象
	encoder:=gob.NewEncoder(&result)
	if err:=encoder.Encode(tx); err!=nil{
		log.Panicf("tx Hash encoded failed %v\n",err)
	}

	//生成哈希值
	hash:=sha256.Sum256(result.Bytes())
	tx.TxHash=hash[:]
}
//生成普通转账交易
func NewSimpleTransaction(from string,to string,amount int,bc*BlockChain) *Transaction {
	var txInputs []*TxInput    //输入列表
	var txOutputs []*TxOutput    //输出列表
	//调用可话费UTXO函数
	money,spendableUTXODic:=bc.FindSpendableUTXO(from,amount)
	fmt.Printf("money : %v\n",money)
	//输入
	for txHash,indexArray:=range spendableUTXODic{
		txHashBytes,err:=hex.DecodeString(txHash)
		if nil!=err{
			log.Panicf("decode string to []byte failed! %v\n",err)
		}
		//遍历索引列表
		for _,index:=range indexArray{
			txInput:=&TxInput{txHashBytes,index,from}
			txInputs=append(txInputs,txInput)
		}
	}
	/*txInput:=&TxInput{[]byte("6f92e00377b91ddb1660fe222eb71" +
		"d24caced56a737a986bf807e6cd885e22f7"),0,from,}
	txInputs=append(txInputs,txInput)*/
	//输出(转账源)
	txOutput:=&TxOutput{amount,to}
	txOutputs=append(txOutputs,txOutput)
	//输出(找零)
	if money>amount{
		txOutput=&TxOutput{money-amount,from}
		txOutputs=append(txOutputs,txOutput)
	}else{
		log.Panicf("余额不足...\n")
	}
	tx:=Transaction{nil,txInputs,txOutputs}
	tx.HashTransaction()
	return &tx
}
//判断指定的交易是否是一个coinbase交易
func (tx *Transaction)IsCoinbaseTransaction() bool {
	return tx.Vins[0].Vout==-1&&len(tx.Vins[0].TxHash)==0
}