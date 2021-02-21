package BLC
//区块链管理文件

//区块链的基本结构
type BlockChain struct{
	Blocks []*Block //区块的切片
}

//初始化区块
func CreateBlockChainWithGenesisBlock()*BlockChain  {
	//生成创世区块
	block:=CreateGenesisBlock([]byte("init blockchain"))
	return &BlockChain{[]*Block{block}}
}
//添加区块
func (bc *BlockChain) AddBlock(height int64,preBlockHash []byte,data []byte)  {
	var newBlock *Block
	newBlock=NewBlock(height,preBlockHash,data)
	bc.Blocks=append(bc.Blocks,newBlock)
}