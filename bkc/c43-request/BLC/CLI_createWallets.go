package BLC

import "fmt"

//创建钱包集合
func (cli *CLI) CreateWallets(nodeID string)  {
	wallets:=NewWallets(nodeID)  //创建一个集合对象
	wallets.CreateWallet(nodeID)
	fmt.Printf("wallets : %v\n",wallets)
}
