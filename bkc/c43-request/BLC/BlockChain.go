package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
)

//区块链管理文件
//数据库名称
const dbName="block_%s.db"
//表名称
const blockTableName="blocks"
//区块链的基本结构
type BlockChain struct{
	//Blocks []*Block //区块的切片
	DB    *bolt.DB    //数据库对象

	Tip   []byte      //保存最新区块的哈希值
}

//判断数据库文件是否存在
func dbExist(nodeID string) bool {
	//生成不同节点的数据库文件
	dbName:=fmt.Sprintf(dbName,nodeID)
	if _,err:=os.Stat(dbName);os.IsNotExist(err){
		//数据库文件不存在
		return false
	}
	return true
}


//初始化区块链
func CreateBlockChainWithGenesisBlock(address string,nodeID string)*BlockChain  {
	if dbExist(nodeID){
		//文件已存在，说明创世区块已存在
		fmt.Println("创世区块已存在...")
		os.Exit(1)
	}
	//保存最新区块的哈希值
	var blockHash []byte
	//1.创建或者打开一个数据库
	//w r x
	dbName:=fmt.Sprintf(dbName,nodeID)
	db,err:=bolt.Open(dbName,0600,nil)

	if err!=nil{
		log.Panicf("create [%s] failed %v\n",dbName,err)
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
/*func (bc *BlockChain) AddBlock(txs []*Transaction)  {
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
}*/
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
				fmt.Printf("\t\t\tvin-PublicKey : %x\n",vin.PublicKey)
				fmt.Printf("\t\t\tvin-Signature : %x\n",vin.Signature)
			}
			fmt.Printf("\t\t输出...\n")
			for _,vout:=range tx.Vouts{
				fmt.Printf("\t\t\tvout-value : %d\n",vout.Value)
				fmt.Printf("\t\t\tvout-Ripems160Hash : %x\n",vout.Ripemd160Hash)
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
func BlockChainObject(nodeID string) *BlockChain {
	//获取DB
	dbName:=fmt.Sprintf(dbName,nodeID)
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
func (blockchain *BlockChain) MineNewBlock(from,to,amount []string,nodeID string)  {
	//搁置交易生成步骤
	var block *Block
	var txs []*Transaction
	//遍历交易的参与者
	for index,address:=range from{
		value,_:=strconv.Atoi(amount[index])
		//生成新的交易
		tx:=NewSimpleTransaction(address,to[index],value,blockchain,txs,nodeID)
		//签名

		//追加到txs的交易列表中去
		txs=append(txs,tx)
		//给予交易发起者（矿工）一定的奖励
		tx=NewCoinbaseTransaction(address)
		txs=append(txs,tx)
	}
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
	//在此处进行交易签名的验证
	//对txs中的每一笔交易的签名都进行验证
	for _,tx:=range txs{
		//验证签名，只要有一笔签名的验证失败，panic
		if blockchain.VerifyTransaction(tx)==false{
			log.Panicf("ERROR : tx [%x] verify failed!\n",tx)
		}
	}

	//通过数据库中最新的区块生成新的区块（交易的打包）
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
//获取指定地址所有已花费输出
func (blockchain *BlockChain) SpentOutput(address string) map[string][]int {
	//已花费输出缓存
	spentTXOutputs:=make(map[string][]int)
	//获取迭代器对象
	bcit:=blockchain.Iterator()
	for{
		block:=bcit.Next()
		for _,tx:=range block.Txs{
			//排除coinbase交易
			if !tx.IsCoinbaseTransaction(){
				for _,in:=range tx.Vins{
					if in.UnLockRipemd160Hash(StringToHash160(address)){
						key:=hex.EncodeToString(in.TxHash)
						//添加到已花费输出的缓存中
						spentTXOutputs[key]=append(spentTXOutputs[key],in.Vout)
					}
					/*if in.CheckPubkeyWithAddress(address){
						key:=hex.EncodeToString(in.TxHash)
						//添加到已花费输出的缓存中
						spentTXOutputs[key]=append(spentTXOutputs[key],in.Vout)
					}*/
				}
			}
		}
		//退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if hashInt.Cmp(big.NewInt(0))==0{
			break
		}
	}
	return spentTXOutputs
}


//查找指定地址的UTXO
/*遍历查找区块链数据库中的每一个区块中的每一个交易
  查找每一个交易中的每一个输出
  判断每个输出是否满足下列条件
  1.属于传入的地址
  2.是否未被花费
    1.首先，遍历一次区块链数据库，将所有已花费的OUTPUT存入一个缓存
    2.再次遍历区块链数据库，检查每一个VOUT是否包含在前面的已花费输出的缓存中
 */


func (blockchain *BlockChain) UnUTXOS(address string,txs []*Transaction) []*UTXO {
	//1.遍历数据库，查找所有与address相关的对象
	//获取迭代器
	//当前地址的未花费输出列表
	var unUTXOS []*UTXO
	bcit:=blockchain.Iterator()
	//获取指定地址所有已花费输出
	spentTXOutputs:=blockchain.SpentOutput(address)
	//缓存迭代
	//查找缓存中的已花费输出
	for _,tx:=range txs{
		//判断coinbaseTransaction
		if !tx.IsCoinbaseTransaction(){
			for _,in:=range tx.Vins{
				//判断用户
				if in.UnLockRipemd160Hash(StringToHash160(address)){
					//添加到已花费输出的map中
					key:=hex.EncodeToString(in.TxHash)
					spentTXOutputs[key]=append(spentTXOutputs[key],in.Vout)
				}


				/*if in.CheckPubkeyWithAddress(address){
					//添加到已花费输出的map中
					key:=hex.EncodeToString(in.TxHash)
					spentTXOutputs[key]=append(spentTXOutputs[key],in.Vout)
				}*/
			}
		}
	}
	//遍历缓存中的UTXO
	for _,tx:=range txs{
		//添加一个缓存输出的跳转
		WorkCacheTx:
		for index,vout:=range tx.Vouts{
			if vout.UnLockScriptPubkeyWithAddress(address){
			//if vout.CheckPubkeyWithAddress(address){
				if len(spentTXOutputs)!=0{
					var isUtxoTx bool  //判断交易是否被其他交易引用
					for txHash,indexArray:=range spentTXOutputs{
						txHashStr:=hex.EncodeToString(tx.TxHash)
						if txHash==txHashStr{
							//当前遍历到交易已经有输出被其他交易的输入所引用
							isUtxoTx=true
							//添加状态变量，判断指定的output是否被引用
							var isSpentUTXO bool
							for _,voutIndex:=range indexArray{
								if index==voutIndex{
									//该输出被引用
									isSpentUTXO=true
									//跳出当前vout判断逻辑，进行下一个输出判断
									continue WorkCacheTx
								}
							}
							if isSpentUTXO==false{
								utxo:=&UTXO{tx.TxHash,index,vout}
								unUTXOS=append(unUTXOS,utxo)
							}
						}
					}
					if isUtxoTx==false{
						//说明当前交易中所有与address相关的outputs都是UTXO
						utxo:=&UTXO{tx.TxHash,index,vout}
						unUTXOS=append(unUTXOS,utxo)
					}
				}else{
					utxo:=&UTXO{tx.TxHash,index,vout}
					unUTXOS=append(unUTXOS,utxo)
				}
			}
		}
	}
	//优先遍历缓存中的UTXO，如果余额足够，直接返回，如果不足，再遍历db文件中的UTXO
	//数据库迭代，不断获取下一个区块
	for{
		block:=bcit.Next()
		//遍历区块中的每笔交易
		for _, tx:=range block.Txs{
			//跳转
			work:
			for index,vout :=range tx.Vouts{
				//index：当前输出在当前交易中的索引位置
				//vout：当前输出
				if vout.UnLockScriptPubkeyWithAddress(address){
				//if vout.CheckPubkeyWithAddress(address){
					//当前vout属于传入地址
					if len(spentTXOutputs)!=0{
						var isSpentOutput bool //默认false
						for txHash,indexArray:=range spentTXOutputs{
							for _,i:=range indexArray{
								//txHash：当前输出所引用的交易哈希
								//indexArray：哈希关联的vout索引列表
								if txHash==hex.EncodeToString(tx.TxHash) && index==i{
									//txHash==hex.EncodeToString(tx.TxHash),
									//说明当前的交易tx至少已经由输出被其他交易的输入引用
									//index==i说明正好是当前的输出被其他交易引用
									//跳转到最外层循环，判断下一个VOUT
									isSpentOutput=true
									continue work
								}
							}
						}
						if isSpentOutput==false{
							utxo:=&UTXO{tx.TxHash,index,vout}
							unUTXOS=append(unUTXOS, utxo)
						}
					}else{
						//将当前所有输出都添加到未花费输出中
						utxo:=&UTXO{tx.TxHash,index,vout}
						unUTXOS=append(unUTXOS,utxo)
					}
				}
			}
		}

		//退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if hashInt.Cmp(big.NewInt(0))==0{
			break
		}
	}
	return unUTXOS
}
//查询余额
func (blockchain *BlockChain) getBalance(address string) int {
	var amount int     //余额
	utxos:=blockchain.UnUTXOS(address,[]*Transaction{})
	for _,utxo:=range utxos{
		amount+=utxo.Output.Value
	}
	return amount
}

//查找指定地址的可用UTXO，超过amount就中断查找
//更新当前数据库中指定地址的UTXO数量
//txs:缓存中的交易列表（用于多笔交易处理）
func (blockchain *BlockChain) FindSpendableUTXO(from string,amount int,txs []*Transaction) (int, map[string][]int) {
	spendableUTXO:=make(map[string][]int)
	var value int
	utxos:=blockchain.UnUTXOS(from,txs)
	//遍历UTXO
	for _,utxo:=range utxos{
		value+=utxo.Output.Value
		//计算交易哈希
		hash:=hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash]=append(spendableUTXO[hash],utxo.Index)
		if value>=amount{
			break
		}
	}

	//所有的都遍历完成，仍然小于amount
	//资金不足
	if value < amount{
		fmt.Printf("地址 [%s] 余额不足，当前余额 [%d]，转账金额 [%d]\n",from,value,amount)
		os.Exit(1)
	}
	return value,spendableUTXO
}

//通过指定的交易哈希查找交易
func (blockchain *BlockChain) FindTransaction(ID []byte) Transaction {
	bcit:=blockchain.Iterator()
	for{
		block:=bcit.Next()
		for _,tx:=range block.Txs{
			if bytes.Compare(ID,tx.TxHash)==0{
				//找到该交易
				return *tx
			}
		}
		//退出
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if big.NewInt(0).Cmp(&hashInt)==0{
			break
		}
	}
	fmt.Printf("没找到交易[%x]\n",ID)
	return Transaction{}
}

//交易签名
func (blockchain *BlockChain) SignTransaction(tx *Transaction,privKey ecdsa.PrivateKey)  {
	//coinbase交易不需要签名
	if tx.IsCoinbaseTransaction(){
		return
	}
	//处理交易的input，查找input所引用的vout所属交易（查找发送者）
	//对我们所花费的每一笔UTXO进行签名
	//存储引用的交易
	prevTxs:=make(map[string]Transaction)
	for _,vin:=range tx.Vins{
		//查找当前交易输入所引用的交易
		tx:=blockchain.FindTransaction(vin.TxHash)
		prevTxs[hex.EncodeToString(tx.TxHash)]=tx

	}
	//签名
	tx.Sign(privKey,prevTxs)
}

//验证签名
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbaseTransaction(){
		return true
	}
	prevTxs:=make(map[string]Transaction)
	//查找输入引用的交易
	for _,vin:=range tx.Vins{
		tx:=bc.FindTransaction(vin.TxHash)
		prevTxs[hex.EncodeToString(tx.TxHash)]=tx
	}
	return tx.Verify(prevTxs)
}
//退出条件
func isBreakLoop(prevBlockHash []byte) bool {
	var hashInt big.Int
	hashInt.SetBytes(prevBlockHash)
	if hashInt.Cmp(big.NewInt(0))==0{
		return true
	}
	return false
}

//查找整条区块链所有已花费输出
func (blockchain *BlockChain) FindAllSpentOutputs() map[string][]*TxInput {
	bcit:=blockchain.Iterator()
	//存储已花费输出
	spentTXOutputs:=make(map[string][]*TxInput)
	for {
		block:=bcit.Next()
		for _,tx:=range block.Txs{
			if !tx.IsCoinbaseTransaction(){
				for _,txInput:=range tx.Vins{
					txHash:=hex.EncodeToString(txInput.TxHash)
					spentTXOutputs[txHash]=append(spentTXOutputs[txHash],txInput)
				}
			}
		}
		if isBreakLoop(block.PreBlockHash){
			break
		}
	}
	return spentTXOutputs
}

//查找整条区块链中所有地址的UTXO
func (blockchain *BlockChain) FindUTXOMap() map[string] *TXOutputs {
	//遍历区块链
	bcit:=blockchain.Iterator()
	//输出集合
	utxoMaps:=make(map[string] *TXOutputs)
	//查找已经花费输出
	spentTXOutputs:=blockchain.FindAllSpentOutputs()
	for{
		block:=bcit.Next()
		for _,tx:=range block.Txs{
			txOutputs:=&TXOutputs{[]*TxOutput{}}
			txHash:=hex.EncodeToString(tx.TxHash)
			//获取每笔交易的vouts
			WorkOutLoop:
			for index,vout:=range tx.Vouts{
				//获取指定交易
				txInputs:=spentTXOutputs[txHash]
				if len(txInputs)>0{
					isSpent:=false
					for _,in:=range txInputs{
						//查找指定输出的所有者
						outPubKey:=vout.Ripemd160Hash
						inPubKey:=in.PublicKey
						if bytes.Compare(outPubKey,Ripemd160Hash(inPubKey))==0{
							if index==in.Vout{
								isSpent=true
								continue WorkOutLoop
							}
						}
					}
					if isSpent==false{
						//当前输出没有被包含的到txInputs中
						txOutputs.TXOutputs=append(txOutputs.TXOutputs,vout)
					}


				}else{
					//没有input引用该交易的输出，则代表当前交易中的所有的输出都是UTXO
					txOutputs.TXOutputs=append(txOutputs.TXOutputs,vout)
				}
			}
			utxoMaps[txHash]=txOutputs
		}
		if isBreakLoop(block.PreBlockHash){
			break
		}
	}
	return utxoMaps
}

//获取当前区块的区块高度
func (bc *BlockChain) GetHeight() int64 {
	return bc.Iterator().Next().Height
}

//获取区块链所有的区块哈希
func (bc *BlockChain) GetBlockHases() [][]byte {
	var blockHashes [][]byte
	bcit:=bc.Iterator()
	for{
		block:=bcit.Next()
		blockHashes=append(blockHashes,block.Hash)
		if isBreakLoop(block.PreBlockHash){
			break
		}
	}
	return blockHashes
}

//获取指定哈希的区块数据
func (bc *BlockChain) GetBlock(hash []byte) []byte {
	var blockByte []byte
	bc.DB.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockTableName))
		if nil!=b{
			blockByte=b.Get(hash)
		}
		return nil
	})
	return blockByte
}

//添加区块
func (bc *BlockChain) AddBlock(block *Block)  {
	err:=bc.DB.Update(func(tx *bolt.Tx) error {
		//1.获取数据表
		b:=tx.Bucket([]byte(blockTableName))
		if nil!=b{
			//判断需要传入的区块是否已经存在
			if b.Get(block.Hash)!=nil{
				//已经存在，不需要添加
				return nil
			}
			//不存在，添加到数据库中
			err:=b.Put(block.Hash,block.Serialize())
			if nil!=err{
				log.Panicf("sync the block failed! %v\n",err)
			}
			blockHash:=b.Get([]byte("1"))
			latestBlock:=b.Get(blockHash)
			rawBlock:=DeserializeBlock(latestBlock)
			if rawBlock.Height<block.Height{
				b.Put([]byte("1"),block.Hash)
				bc.Tip=block.Hash
			}
		}
		return nil
	})
	if nil!=err{
		log.Panicf("update the db when insert the new block failed! %v\n",err)
	}
	fmt.Println("the new block is added!")
}