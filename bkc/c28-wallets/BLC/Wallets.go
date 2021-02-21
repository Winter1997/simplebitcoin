package BLC

//钱包集合管理的文件

//实现钱包集合的基本结构
type Wallets struct{
	//key : 地址
	//value : 钱包结构
	Wallets map[string] *Wallet
}


//初始化钱包集合
func NewWallets() *Wallets {
	wallets:=&Wallets{}
	wallets.Wallets=make(map[string] *Wallet)
	return wallets
}

//添加新的钱包到集合中
func (wallets *Wallets) CreateWallet()  {
	//1.创建钱包
	wallet:=NewWallet()
	//2.添加
	wallets.Wallets[string(wallet.GetAddress())]=wallet
}