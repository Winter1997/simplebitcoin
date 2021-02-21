package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
)

//钱包管理相关文件

//钱包基本结构
type Wallet struct{
	//1.私钥
	PrivateKey ecdsa.PrivateKey
	//2.公钥
	PublicKey  []byte

}

//创建一个钱包
func NewWallet() *Wallet  {
	//公钥-私钥赋值
	privateKey,publicKey:=newKeyPair()
	return &Wallet{PrivateKey: privateKey,PublicKey: publicKey}
}
//通过钱包生成公钥-私钥对
func newKeyPair()(ecdsa.PrivateKey,[]byte)  {
	//1.获取一个椭圆
	curve:=elliptic.P256()
	//2.通过椭圆相关算法生成私钥
	priv,err:=ecdsa.GenerateKey(curve,rand.Reader)
	if nil!=err{
		log.Panicf("ecdsa generate private key failed! %v\n",err)
	}
	//3.通过私钥生成公钥
	pubKey:=append(priv.PublicKey.X.Bytes(),priv.PublicKey.Y.Bytes()...)
	return *priv,pubKey
}