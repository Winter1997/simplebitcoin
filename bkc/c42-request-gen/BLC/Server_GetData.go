package BLC

//请求指定区块
type GetData struct {
	AddrFrom    string    //从哪一个地址请求
	ID          []byte    //区块哈希
}
