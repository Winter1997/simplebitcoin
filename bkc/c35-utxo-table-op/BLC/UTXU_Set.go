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

func (txOutputs *TXOutputs) Serialize() []byte {
	var result bytes.Buffer
	encoder:=gob.NewEncoder(&result)
	if err:=encoder.Encode(txOutputs);nil!=err{
		log.Panicf("seralize the utxo failed! %v\n",err)
	}
	return result.Bytes()
}


//更新


//查找


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