package main

import (
	"encoding/json"
	"fmt"

	"github.com/Stan/168proj/cakecoin/utils"
)

func main() {

	fakeNet := NewFakeNet()
	Alice := NewClient("Alice", fakeNet, nil, nil)
	Bob := NewClient("Bob", fakeNet, nil, nil)
	bc := BlockChain{}
	clientBalanceMap := make(map[*Client]int)
	clientBalanceMap[Alice] = 200
	clientBalanceMap[Bob] = 100
	g := bc.MakeGenesis(5, 25, 5, 6, clientBalanceMap, nil)
	fmt.Println(g.Serialize())
	fmt.Println(Alice.confirmedBalance())
	fmt.Println(Alice.availableGold())
	outputs := []Output{Output{3.0, "randomstring"}, Output{4.0, "randomstring2"}}
	Alice.postTransaction(outputs, 5)
	fmt.Println(Alice.availableGold())
	alStr, _ := json.Marshal(Alice)
	fmt.Println(string(alStr))
	bobStr, _ := json.Marshal(Bob)
	fmt.Println(string(bobStr))

	privKey := utils.GenerateKeypair()
	addr1 := utils.CalculateAddress(&privKey.PublicKey)
	newBlock := bc.MakeBlock(addr1, g, nil, nil)
	for !newBlock.hasValidProof() {
		newBlock.Proof++
	}
	fmt.Println(string(g.id()))
	fmt.Println(string(newBlock.Serialize()))

	fakeNet.register(Alice, Bob)
	fakeNet.broadcast(PROOF_FOUND, newBlock)

	aStr, err := json.Marshal(Alice)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(aStr))
	fmt.Println("done")

	/**
	emitter := emission.NewEmitter()
	hello := func(b *Block) {
		fmt.Println(b.Serialize())
	}

	emitter.On(PROOF_FOUND, hello).
		Emit(PROOF_FOUND, newBlock)


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

	fmt.Println("---------------Transactions---------------")
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
	fmt.Println(string(trans.Id()))
	fmt.Println("")

	fmt.Println("---------------Block---------------")
	b := NewBlock(addr1, nil, big.NewInt(5), 25)
	b.ChainLength = 1
	fmt.Println(b.Serialize())
	deserializeBlock([]byte(b.Serialize()))
	fmt.Println(b.IsGenesisBlock())
	fmt.Println(b.contains(*trans))
	fmt.Println("")

	b.Balances[addr1] = 25
	b.addTransaction(trans, 0)
	for id, _ := range b.Transactions {
		fmt.Println([]byte(id))
		fmt.Println(trans.Id())
		fmt.Println(bytes.Compare([]byte(id), trans.Id()))
	}
	fmt.Println("---------------------")
	fmt.Println(b.Serialize())
	db := deserializeBlock([]byte(b.Serialize()))
	fmt.Println(db.Balances)
	fmt.Println(db.Transactions)
	for id, _ := range db.Transactions {
		fmt.Println([]byte(id))
		fmt.Println(trans.Id())
		fmt.Println(bytes.Compare([]byte(id), trans.Id()))
	}
	transStr2, err2 := json.Marshal(db.Transactions[string(trans.Id())])
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println(string(transStr2))
	fmt.Println(db.contains(*trans))

	balanceMap := make(map[string]int)
	balanceMap[addr1] = 25
	bc := BlockChain{}
	gb := bc.MakeGenesis(3, 25, 5, 6, nil, balanceMap)
	gbString, err3 := json.Marshal(gb)
	if err3 != nil {
		fmt.Println(err3)
	}
	fmt.Println(string(gbString))
	**/
}
