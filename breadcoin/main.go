package main

import (
	"encoding/json"
	"fmt"

	"github.com/Stan/168proj/breadcoin/utils"
)

func main() {
	text := "I am lending you $100 for 10 years"
	h := utils.Hash(text)
	fmt.Printf("%x\n", h)

	privKey := utils.GenerateKeypair()
	addr1 := utils.CalculateAddress(&privKey.PublicKey)
	fmt.Println(addr1)
	fmt.Println(utils.AddressMatchesKey(addr1, &privKey.PublicKey))

	signature := utils.Sign(privKey, text)
	fmt.Println(signature)
	fmt.Println(utils.VerifySignature(&privKey.PublicKey, text, signature))

	outputs := []Output{Output{3.0, "randomstring"}}
	trans := NewTransaction(addr1, 0, &privKey.PublicKey, nil, 5.0, outputs)
	transStr, err := json.Marshal(trans)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(transStr))
	fmt.Println(trans.TotalOutput())
	fmt.Println(trans.ValidSignature())
	trans.Sign(privKey)
	fmt.Println(trans.ValidSignature())
	fmt.Println(trans.Id())
}
