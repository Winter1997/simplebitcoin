package BLC

//交易的输出管理

//输出结构
type TxOutput struct{
	//金额
	value   int
	//用户名（UTXO）的所有者
	ScriptPubkey  string

}
