package main

import (
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/ripemd160"
)

func main()  {
	//sha256
	hash:=sha256.New()
	hash.Write([]byte("eth1804"))
	bytes:=hash.Sum(nil)
	fmt.Printf("sha256 : %x\n",bytes)

	//ripemd160
	var r160 = ripemd160.New()
	r160.Write(bytes)
	bytesRipemd:=r160.Sum(nil)
	fmt.Printf("ripemd160 : %x\n",bytesRipemd)
}
