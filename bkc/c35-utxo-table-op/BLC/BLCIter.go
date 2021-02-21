package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

//区块链迭代器管理文件

//迭代器基本结构
type BlockChainIterator struct{
	DB  *bolt.DB   //迭代目标
	CurrentHash []byte  //当前接待目标的哈希
}

//创建迭代器对象
func (blc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{DB: blc.DB,CurrentHash: blc.Tip}
}
//实现迭代函数next，获取到每一个区块
func (bcit *BlockChainIterator) Next() *Block {
	var block *Block

	err:=bcit.DB.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockTableName))
		if nil!=b{
			currentBlockBytes:=b.Get(bcit.CurrentHash)
			block=DeserializeBlock(currentBlockBytes)
			//更新迭代器中区块的哈希值
			bcit.CurrentHash=block.PreBlockHash
		}
		return nil
	})
	if nil!=err{
		log.Panicf("iterator the db failed! %v\n",err)
	}
	return block
}
