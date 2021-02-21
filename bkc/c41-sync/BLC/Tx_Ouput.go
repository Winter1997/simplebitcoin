package BLC

import "bytes"

//交易的输出管理

//输出结构
type TxOutput struct{
	//金额
	Value   int

	//ScriptPubkey  string
	//用户名（UTXO）的所有者
	Ripemd160Hash  []byte
}
//验证当前UTXO是否属于指定的地址
/*func (txOutput *TxOutput)CheckPubkeyWithAddress(address string) bool {
	return address==txOutput.ScriptPubkey
}*/

//output身份验证
func (TxOutput *TxOutput) UnLockScriptPubkeyWithAddress(address string ) bool {
	//转换
	hash160:=StringToHash160(address)
	return bytes.Compare(hash160,TxOutput.Ripemd160Hash)==0
}

//新建output对象
func NewTxOutput(value int,address string) *TxOutput {
	txOutput:=&TxOutput{}
	hash160:=StringToHash160(address)
	txOutput.Value=value
	txOutput.Ripemd160Hash=hash160
	return txOutput
}

