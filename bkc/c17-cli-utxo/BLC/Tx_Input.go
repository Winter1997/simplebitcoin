package BLC

//交易输入管理

//输入结构
type TxInput struct{
	//交易哈希（不是指当前交易的哈希）
	TxHash  []byte
	//引用的上一笔交易的输出索引号
	Vout    int
	//用户名
	ScriptSig  string
}
//验证引用的地址是否匹配
func (txInput *TxInput)CheckPutkeyWithAddress(address string) bool {
	return address==txInput.ScriptSig
}