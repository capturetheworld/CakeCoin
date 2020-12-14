package main

import (
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"

	//"./utils"
	"github.com/Stan/168proj/cakecoin/utils"
)

type Output struct {
	Amount int
	Addr   string
}

type Transaction struct {
	From    string
	Nonce   int
	PubKey  *rsa.PublicKey
	Sig     []byte
	Fee     int
	Outputs []Output
	//data idk what to do for this for now
}

func (t Transaction) Id() []byte {
	transStr := t.toJson()

	return []byte(hex.EncodeToString(utils.Hash("TX" + string(transStr))))
}

func (t Transaction) toJson() []byte {
	var tmp struct {
		From    string
		Nonce   int
		PubKey  *rsa.PublicKey
		Fee     int
		Outputs []Output
	}
	tmp.From = t.From
	tmp.Nonce = t.Nonce
	tmp.PubKey = t.PubKey
	tmp.Fee = t.Fee
	tmp.Outputs = t.Outputs
	val, err := json.Marshal(&tmp)
	if err != nil {
		panic(err)
	}
	return val
}

func (t *Transaction) Sign(privKey *rsa.PrivateKey) {
	t.Sig = utils.Sign(privKey, string(t.Id()))
}

func (t Transaction) ValidSignature() bool {
	if t.Sig != nil {
		if !utils.AddressMatchesKey(t.From, t.PubKey) {
			fmt.Println("doesn't match key")
			return false
		} else if !utils.VerifySignature(t.PubKey, string(t.Id()), t.Sig) {
			fmt.Println("invalid sig")
			return false
		}
	}
	return true
	//return t.Sig != nil && utils.AddressMatchesKey(t.From, t.PubKey) && utils.VerifySignature(t.PubKey, string(t.Id()), t.Sig)
}

func (t Transaction) SufficientFunds(b Block) bool {
	return t.TotalOutput() <= b.Balances[t.From]
}

func (t Transaction) TotalOutput() int {
	total := t.Fee
	for _, output := range t.Outputs {
		total += output.Amount
	}
	return total
}

//NEED DATA PLZZZZ ADDDDDD YOOOOOO
func NewTransaction(from string, nonce int, pubKey *rsa.PublicKey, sig []byte, fee int, outputs []Output) *Transaction {

	//didn't do nil checking
	//fmt.Println("working")
	newTransaction := Transaction{from, nonce, pubKey, sig, fee, outputs}
	return &newTransaction
}
