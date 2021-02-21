package main

import (
	"bitcoin/bkc/c28-wallets/BLC"
	"fmt"
)

func main()  {
	result:=BLC.Base58Encode([]byte("this is the example"))
	fmt.Printf("result : %s\n",result)
	decodeResult:=BLC.Base58Decode([]byte("1nl2SLMErZakmBni8xhSXtimREn"))
	fmt.Printf("decodeResult : %s\n",decodeResult)
}
