package main

import (
	"fmt"

	"github.com/user/168proj/breadcoin/utils"
)

func main() {
	text := "I am lending you $100 for 10 years"
	h := utils.Hash(text)
	fmt.Printf("%x\n", h)

	privKey := utils.GenerateKeypair()
	addr1 := utils.CalculateAddress(privKey)
	fmt.Println(addr1)
	fmt.Println(utils.AddressMatchesKey(addr1, privKey))

	signature := utils.Sign(privKey, text)
	fmt.Println(signature)
	fmt.Println(utils.VerifySignature(privKey, text, signature))

}
