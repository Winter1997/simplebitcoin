package BLC


//展示
type Inv struct {
	AddrFrom    string    //当前节点的地址
	Hashes      [][]byte  //当前节点所拥有的区块的hash列表
}