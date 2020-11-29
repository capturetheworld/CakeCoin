package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"

	"github.com/Stan/168proj/breadcoin/utils"
)

type Output struct {
	Amount float64
	Addr   string
}

type Transaction struct {
	From    string
	Nonce   int
	PubKey  *rsa.PublicKey
	Sig     []byte `json:"-"`
	Fee     float64
	Outputs []Output
	//data idk what to do for this for now
}

func (t Transaction) Id() []byte {
	transStr, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return utils.Hash("TX" + string(transStr))
}

func (t *Transaction) Sign(privKey *rsa.PrivateKey) {
	t.Sig = utils.Sign(privKey, string(t.Id()))
}

func (t Transaction) ValidSignature() bool {
	return t.Sig != nil && utils.AddressMatchesKey(t.From, t.PubKey) && utils.VerifySignature(t.PubKey, string(t.Id()), t.Sig)
}

func (t Transaction) SufficientFunds(b Block) bool {
	return t.TotalOutput() <= b.Balances[t.From]
}

func (t Transaction) TotalOutput() float64 {
	total := t.Fee
	for _, output := range t.Outputs {
		total += output.Amount
	}
	return total
}

//NEED DATA PLZZZZZ ADDDDDD YOOOOOO
func NewTransaction(from string, nonce int, pubKey *rsa.PublicKey, sig []byte, fee float64, outputs []Output) *Transaction {

	//didn't do nil checking
	fmt.Println("working")
	newTransaction := Transaction{from, nonce, pubKey, sig, fee, outputs}
	return &newTransaction
}
