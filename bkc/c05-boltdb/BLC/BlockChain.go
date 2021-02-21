package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

//区块链管理文件
//数据库名称
const dbName="block.db"
//表名称
const blockTableName="blocks"
//区块链的基本结构
type BlockChain struct{
	//Blocks []*Block //区块的切片
	DB    *bolt.DB    //数据库对象

	Tip   []byte      //保存最新区块的哈希值
}

//初始化区块
func CreateBlockChainWithGenesisBlock()*BlockChain  {
	//保存最新区块的哈希值
	var blockHash []byte
	//1.创建或者打开一个数据库
	//w r x
	db,err:=bolt.Open(dbName,0600,nil)

	if err!=nil{
		log.Panic("create [%s] failed %v\n",dbName,err)
	}
	//2.创建桶,把生成的区块放入到数据库中
	db.Update(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockTableName))
		if b==nil{
			//没找到桶
			b,err:=tx.CreateBucket([]byte(blockTableName))
			if err!=nil{
				log.Panicf("create bucket [%s] failed %v\n",blockTableName,err)
			}
			//生成创世区块
			genesiBlock:=CreateGenesisBlock([]byte("init blockchain"))
			//存储
			//1.key,value分别以什么数据代表--hash
			//2.如何把block结构存入到数据库中--序列化
			err1:=b.Put(genesiBlock.Hash,genesiBlock.Serialize())
			if err1!=nil{
				log.Panicf("insert genesi block failed %v\n",err1)
			}
			blockHash=genesiBlock.Hash
			//存储最新区块的哈希
			err2:=b.Put([]byte("1"),genesiBlock.Hash)
			if err2!=nil{
				log.Panicf("save the latest hash of genesis block %v\n",err2)
			}
		}
		return nil
	})
	//3.把创世区块存入到数据库中
	return &BlockChain{DB:db,Tip: blockHash}
}
//添加区块
func (bc *BlockChain) AddBlock(data []byte)  {
	//更新区块数据（insert）
	err2:=bc.DB.Update(func(tx *bolt.Tx) error {
		//1.获取数据库桶
		b:=tx.Bucket([]byte(blockTableName))
		if nil!=b{
			//2.获取最后插入的区块
			blockBytes:=b.Get(bc.Tip)
			//3.区块数据反序列化
			latest_block:=DeserializeBlock(blockBytes)
			//3.新建区块
			newBlock:=NewBlock(latest_block.Height+1,latest_block.Hash,data)
			//4.存入数据库
			err:=b.Put(newBlock.Hash,newBlock.Serialize())
			if err!=nil{
				log.Panicf("insert the new block to db failed %v\n",err)
			}
			//更新最新区块的哈希（数据库）
			err1:=b.Put([]byte("1"),newBlock.Hash)
			if err1!=nil{
				log.Panicf("updata the latest block hash to db failed %v\n",err1)
			}
			//更新区块连对象中的最新区块哈希
			bc.Tip=newBlock.Hash
		}
		return nil
	})
	if nil!=err2{
		log.Panicf("insert block to db failed%v\n",err2)
	}
}