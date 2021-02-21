package BLC

import (
	"fmt"
	"testing"
)

func TestNewWallet(t *testing.T) {
	wallet:=NewWallet()
	fmt.Printf("private key : %v\n",wallet.PrivateKey)
	fmt.Printf("public key : %v\n",wallet.PublicKey)
	fmt.Printf("wallet : %v\n",wallet)
}