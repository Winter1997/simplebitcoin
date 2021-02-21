package BLC

import "fmt"

//查询余额
func (cli *CLI) getBalance(from string,nodeID string)  {
	//查找该地址UTXO
	//获取区块链对象
	blockchain:=BlockChainObject(nodeID)
	defer blockchain.DB.Close() //关闭实例对象
	utxoSet:=UTXOSet{Blockchain: blockchain}
	amount:=utxoSet.GetBalance(from)
	fmt.Printf("\t地址 [%s] 的余额: [%d]\n",from,amount)


}