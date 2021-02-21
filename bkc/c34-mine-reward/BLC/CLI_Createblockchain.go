package BLC
//初始化区块链
func (cli *CLI)createBlockchain(address string) {
	CreateBlockChainWithGenesisBlock(address)
}