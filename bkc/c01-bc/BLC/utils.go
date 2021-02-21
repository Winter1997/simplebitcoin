package BLC

//公共模块

import (
	"bytes"
	"encoding/binary"
	"log"
)

//实现int64转[]byte
func IntToHex(data int64) []byte {
	buffer:=new(bytes.Buffer)
	err:=binary.Write(buffer,binary.BigEndian,data)
	if err!=nil{
		log.Panic("int transact to []byte failed! %v\n",err)
	}
	return buffer.Bytes()
}