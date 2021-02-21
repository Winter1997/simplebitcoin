package BLC

import (
	"bytes"
	"crypto/sha256"
	"time"
)

//区块基本结构与功能管理文件
type Block struct{
	TimeStamp int64  //区块时间戳，代表区块时间
	Hash      []byte //当前区块哈希
	PreBlockHash []byte //前区块哈希
	Height    int64   //区块高度
	Data      []byte  //交易数据
	Nonce     int64   //在运行pow时生成的哈希变化值，也代表pow运行时动态修改的数据
}

//新建区块
func NewBlock(height int64,preBlockHash []byte,data []byte) *Block{
	var block Block
	block=Block{
		TimeStamp: time.Now().Unix(),
		Hash: nil,
		PreBlockHash: preBlockHash,
		Height: height,
		Data: data,
	}
	//生成哈希
	block.SetHash()
	//替换setHash
	//通过POW生成新的哈希值
	pow:=NewProofOfWork(&block)
	//执行工作量证明的算法
	hash,nonce:=pow.Run()
	block.Hash=hash
	block.Nonce=int64(nonce)
	return &block
}
//计算区块
func (b *Block) SetHash(){
	//调用sha256实现哈希生成
	//实现int-hash
	timeStampBytes:=IntToHex(b.TimeStamp)
	heightBytes:=IntToHex(b.Height)
	blockBytes:=bytes.Join([][]byte{
		heightBytes,
		timeStampBytes,
		b.PreBlockHash,
		b.Data,
		b.Hash,
	},[]byte{})
	hash:=sha256.Sum256(blockBytes)
	b.Hash=hash[:]
}

//生成创世区块
func CreateGenesisBlock(data []byte) *Block {
	return NewBlock(1,nil,data)
}