package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
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

//判断数据库文件是否存在
func dbExist() bool {
	if _,err:=os.Stat(dbName);os.IsNotExist(err){
		//数据库文件不存在
		return false
	}
	return true
}


//初始化区块链
func CreateBlockChainWithGenesisBlock(address string)*BlockChain  {
	if dbExist(){
		//文件已存在，说明创世区块已存在
		fmt.Println("创世区块已存在...")
		os.Exit(1)
	}
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
			//生成一个coinbase交易
			txCoinbase:=NewCoinbaseTransaction(address)
			//生成创世区块
			genesiBlock:=CreateGenesisBlock([]*Transaction{txCoinbase})
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
//添加区块到区块链中
func (bc *BlockChain) AddBlock(txs []*Transaction)  {
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
			newBlock:=NewBlock(latest_block.Height+1,latest_block.Hash,txs)
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
//遍历数据库，输出所有区块信息
func (bc *BlockChain) PrintChain()  {
	//读取数据库
	fmt.Println("打印区块链完整信息...")
	var curBlock *Block
	bcit:=bc.Iterator()  //获取迭代器对象
	//循环读取
	//退出条件
	for {
		fmt.Println("\t-------------------------------------------------------------------------------------")
		curBlock=bcit.Next()
		fmt.Printf("\tHash : %x\n",curBlock.Hash)
		fmt.Printf("\tPrevBlockHash : %x\n",curBlock.PreBlockHash)
		fmt.Printf("\tTimeStamp : %v\n",curBlock.TimeStamp)
		fmt.Printf("\tHeight : %d\n",curBlock.Height)
		fmt.Printf("\tNounce : %d\n",curBlock.Nonce)
		fmt.Printf("\tTxs : %v\n",curBlock.Txs)
		for _,tx:=range curBlock.Txs{
			fmt.Printf("\t\ttx-hash : %x\n",tx.TxHash)
			fmt.Printf("\t\t输入...\n")
			for _,vin:=range tx.Vins{
				fmt.Printf("\t\t\tvin-txHash : %x\n",vin.TxHash)
				fmt.Printf("\t\t\tvin-vout : %v\n",vin.Vout)
				fmt.Printf("\t\t\tvin-scriptSig : %s\n",vin.ScriptSig)
			}
			fmt.Printf("\t\t输出...\n")
			for _,vout:=range tx.Vouts{
				fmt.Printf("\t\t\tvout-value : %d\n",vout.Value)
				fmt.Printf("\t\t\tvout-scriptPubkey : %s\n",vout.ScriptPubkey)
			}
		}
		//退出条件
		//转换为big.int
		var hashInt big.Int
		hashInt.SetBytes(curBlock.PreBlockHash)
		//比较
		if big.NewInt(0).Cmp(&hashInt)==0{
			//遍历到创世区块
			break
		}
	}
}
//获取一个blockchain对象
func BlockChainObject() *BlockChain {
	//获取DB
	db,err:=bolt.Open(dbName,0600,nil)
	if nil!=err{
		log.Panicf("open the bd [%s] failed! %v\n",dbName,err)
	}
	//获取Tip
	var tip []byte
	err=db.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockTableName))
		if nil!=b{
			tip=b.Get([]byte("1"))
		}
		return nil
	})
	if nil!=err{
		log.Panicf("get the blockchain object failed! %v\n",err)
	}
	return &BlockChain{DB:db,Tip: tip}
 }
 //实现挖矿功能
 //通过接收交易生成区块
func (blockchain *BlockChain) MineNewBlock(from,to,amount []string)  {
	//搁置交易生成步骤
	var block *Block
	var txs []*Transaction
	value,_:=strconv.Atoi(amount[0])
	//生成新的交易
	tx:=NewSimpleTransaction(from[0],to[0],value)
	//追加到txs的交易列表中去
	txs=append(txs,tx)
	//从数据库中获取最新一个区块
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockTableName))
		if nil!=b{
			hash:=b.Get([]byte("1"))
			//获取最新区块
			blockBytes:=b.Get(hash)
			//反序列化
			block=DeserializeBlock(blockBytes)
		}
		return nil
	})
	//通过数据库中最新的区块生成新的区块
	block=NewBlock(block.Height+1,block.Hash,txs)
	//持久化新生成的区块到数据库中
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockTableName))
		if nil!=b{
			err:=b.Put(block.Hash,block.Serialize())
			if nil!=err{
				log.Panicf("update the new block to db failed! %v\n",err)
			}
			//更新最新区块的哈希值
			err=b.Put([]byte("1"),block.Hash)
			if nil!=err{
				log.Panicf("update the latest block hash to db failed! %v\n",err)
			}
			blockchain.Tip=block.Hash
		}
		return nil
	})
}
//查找指定地址的UTXO
func (blockchain *BlockChain) UnUTXOS(address string) []*TxOutput {
	fmt.Printf("exec the UnUTXOS function\n")
	return nil
}