package BLC

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//UTXO持久化管理
//用于存入utxo的bucket
const utxoTableName="utxoTable"

//utxoSet结构保存（指定区块链中所有的UTXO）

type UTXOSet struct {
	Blockchain *BlockChain
}
//输出集序列化
func (txOutputs *TXOutputs) Serialize() []byte {
	var result bytes.Buffer
	encoder:=gob.NewEncoder(&result)
	if err:=encoder.Encode(txOutputs);nil!=err{
		log.Panicf("seralize the utxo failed! %v\n",err)
	}
	return result.Bytes()
}
//输出集返回序列化
func DeserializeTXOutputs(txOutputsBytes []byte) *TXOutputs {
	var txOutputs TXOutputs
	decoder:=gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	if err:=decoder.Decode(&txOutputs);nil!=err{
		log.Panicf("deserialize the struct utxo failed! %v\n",err)
	}
	return &txOutputs
}

//更新

//查询余额
func (utxoSet *UTXOSet) GetBalance(address string) int {
	UTXOS:=utxoSet.FindUTXOWithAddress(address)
	var amount int
	for _,utxo:=range UTXOS{
		fmt.Printf("utxo-txhash:%x\n",utxo.TxHash)
		fmt.Printf("utxo-index:%x\n",utxo.Index)
		fmt.Printf("utxo-Ripemd160Hash:%x\n",utxo.Output.Ripemd160Hash)
		fmt.Printf("utxo-Value:%x\n",utxo.Output.Value)
		amount+=utxo.Output.Value
	}
	return amount
}
//查找
func (utxoSet *UTXOSet) FindUTXOWithAddress(address string) []*UTXO {
	var utxos []*UTXO
	err:=utxoSet.Blockchain.DB.View(func(tx *bolt.Tx) error {
		//1.获取utxotable表
		b:=tx.Bucket([]byte(utxoTableName))
		if nil!=b{
			//cursor
			c:=b.Cursor()
			//通过游标遍历boltdb数据库中的数据
			for k,v:=c.First();k!=nil;k,v=c.Next(){
				txOutputs:=DeserializeTXOutputs(v)
				for _,utxo:=range txOutputs.TXOutputs{
					if utxo.UnLockScriptPubkeyWithAddress(address){
						utxo_signle:=UTXO{Output: utxo}
						utxos=append(utxos,&utxo_signle)
					}
				}
			}
		}
		return nil
	})
	if nil!=err{
		log.Panicf("find the utxo of [%s] failed! %v\n",address,err)
	}
	return utxos
}

//重置
func (utxoSet *UTXOSet) ResetUTXOSet()  {
	//在第一次创建的时候就更新utxo table
	utxoSet.Blockchain.DB.Update(func(tx *bolt.Tx) error {
		//查找utxo table
		b:=tx.Bucket([]byte(utxoTableName))
		if nil!=b{
			err:=tx.DeleteBucket([]byte(utxoTableName))
			if nil!=err{
				log.Panicf("delete the utxo table failed! %v\n",err)
			}
		}
		//创建
		bucket,err:=tx.CreateBucket([]byte(utxoTableName))
		if nil!=err{
			log.Panicf("create bucket failed! %v\n",err)
		}
		if nil!=bucket{
			//查找当前所有UTXO
			txOutputMap:=utxoSet.Blockchain.FindUTXOMap()
			for keyHash,outputs:=range txOutputMap{
				//将所有UTXO存入
				txHash,_:=hex.DecodeString(keyHash)
				fmt.Printf("KeyHash : %x\n",txHash)
				//存入utxo table
				err:=bucket.Put(txHash,outputs.Serialize())
				if nil!=err{
					log.Panicf("put the utxo into table failed! %v\n",err)
				}
			}
		}
		return nil
	})
}
