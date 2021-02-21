package BLC

//重置utxo table
func (cli *CLI)TestResetUTXO()  {
	blockchain:=BlockChainObject()
	defer blockchain.DB.Close()
	utxoSet:=UTXOSet{Blockchain: blockchain}
	utxoSet.ResetUTXOSet()
}

//重置



//查找
func (cli *CLI)TestFindUTXOMap()  {
	
}