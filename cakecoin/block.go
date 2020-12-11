package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"time"

	"./utils"
)

type Block struct {
	PrevBlockHash  []byte
	Target         *big.Int `json:"-"`
	Balances       map[string]int
	NextNonce      map[string]int `json:"-"`
	Transactions   map[string]*Transaction
	ChainLength    int
	Timestamp      time.Time
	RewardAddr     string
	CoinbaseReward int `json:"-"`
	Proof          int
}

func (b Block) hashVal() []byte {
	return []byte(hex.EncodeToString(utils.Hash(b.Serialize())))
}

func (b Block) totalRewards() int {
	total := int(b.CoinbaseReward)
	for _, output := range b.Transactions {
		total += output.Fee
	}
	return total
}

func (b Block) balanceOf(addr string) int {
	return b.Balances[addr]
}

func (b Block) IsGenesisBlock() bool {
	return b.ChainLength == 0
}

func (b Block) hasValidProof() bool {
	h := string(b.hashVal())
	n := new(big.Int)

	n, err := n.SetString(h, 16)
	if !err {
		fmt.Println("Could not set hash of block to big Int")
		return false
	}

	//fmt.Println(b.Proof)
	//fmt.Println(n)
	//fmt.Println(b.Target)
	return n.Cmp(b.Target) < 0
}

func (b Block) Serialize() string {
	/**
	jsonByte, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	**/
	return string(b.toJson())
}

func (b Block) toJson() []byte {
	var jsonFile []byte

	size := reflect.ValueOf(b).NumField()

	jsonFile = append(jsonFile, '{')

	includeFields := make(map[string]bool)
	includeFields["ChainLength"] = true
	includeFields["Timestamp"] = true
	if b.IsGenesisBlock() {
		includeFields["Balances"] = true
	} else {
		includeFields["Transactions"] = true
		includeFields["PrevBlockHash"] = true
		includeFields["Proof"] = true
		includeFields["RewardAddr"] = true
	}

	for i := 0; i < size; i++ {
		structValue := reflect.ValueOf(b).Field(i)
		var fieldName string = reflect.TypeOf(b).Field(i).Name
		if marshalledField, err := json.Marshal((structValue).Interface()); err != nil {
			panic("CAN'T CONVERT BLOCK TO JSON")
		} else {
			if includeFields[fieldName] {
				jsonFile = append(jsonFile, '"')
				jsonFile = append(jsonFile, []byte(fieldName)...)
				jsonFile = append(jsonFile, '"')
				jsonFile = append(jsonFile, ':')
				jsonFile = append(jsonFile, (marshalledField)...)
				if i+1 != len(includeFields) {
					jsonFile = append(jsonFile, ',')
				}
			}
		}
	}
	jsonFile = append(jsonFile, '}')
	return jsonFile
}

func (b Block) id() []byte {
	return b.hashVal()
}

func (b Block) contains(tx Transaction) bool {
	_, found := b.Transactions[string(tx.Id())]
	if found {
		return true
	}
	return false
}

//ONCE CLIENT IS DONE CHANGE THIS!!!!!!!!!!!!!!!!!
func (b *Block) addTransaction(tx *Transaction, client int) bool {
	if _, found := b.Transactions[string(tx.Id())]; found {
		fmt.Println(string(tx.Id()) + " is a duplicate")
		return false
	} else if tx.Sig == nil {
		fmt.Println(string(tx.Id()) + " is unsigned")
		return false
	} else if !tx.ValidSignature() {
		fmt.Println(string(tx.Id()) + " has an invalid signature")
		return false
	} else if !tx.SufficientFunds(*b) {
		fmt.Println(string(tx.Id()) + " not enough gold for this transactions")
		return false
	}
	nonce := 0
	if val, found := b.NextNonce[string(tx.Id())]; found {
		nonce = val
	}

	if tx.Nonce < nonce {
		fmt.Println(string(tx.Id()) + " is replayed")
		return false
	} else if tx.Nonce > nonce {
		fmt.Println(string(tx.Id()) + " out of order")
		return false
	} else {
		b.NextNonce[tx.From] = nonce + 1
	}

	s := string(tx.Id())
	//fmt.Println(tx.Id())
	b.Transactions[s] = tx
	//fmt.Println(b.Transactions)

	senderBalance := b.balanceOf(tx.From)
	b.Balances[tx.From] = senderBalance - tx.TotalOutput()

	for _, output := range tx.Outputs {
		oldBalance := b.balanceOf(output.Addr)
		b.Balances[output.Addr] = output.Amount + oldBalance
	}

	return true
}

func (b *Block) rerun(prevBlock *Block) bool {
	b.Balances = make(map[string]int)
	b.NextNonce = make(map[string]int)
	txs := b.Transactions
	b.Transactions = make(map[string]*Transaction)
	for index, element := range prevBlock.Balances {
		b.Balances[index] = element
	}
	for index, element := range prevBlock.NextNonce {
		b.NextNonce[index] = element
	}

	winnerBalance := b.balanceOf(prevBlock.RewardAddr)
	if prevBlock.RewardAddr != "" {
		b.Balances[prevBlock.RewardAddr] = winnerBalance + prevBlock.totalRewards()
	}

	for _, element := range txs {
		success := b.addTransaction(element, 0)
		if !success {
			return false
		}
	}

	return true
}

func NewBlock(rewardAddr string, prevBlock *Block, target *big.Int, coinbaseReward int) *Block {
	var prevBlockHash []byte = nil
	balances := make(map[string]int)
	nextNonce := make(map[string]int)
	transactions := make(map[string]*Transaction)
	chainLength := 0
	if prevBlock != nil {
		prevBlockHash = prevBlock.hashVal()
		for index, element := range prevBlock.Balances {
			balances[index] = element
		}
		for index, element := range prevBlock.NextNonce {
			nextNonce[index] = element
		}
		chainLength = prevBlock.ChainLength + 1
	}

	newBlock := Block{PrevBlockHash: prevBlockHash, Target: target, Balances: balances, NextNonce: nextNonce, Transactions: transactions, ChainLength: chainLength, Timestamp: time.Now(), RewardAddr: rewardAddr, CoinbaseReward: coinbaseReward}

	if prevBlock != nil && prevBlock.RewardAddr != "" {
		winnerBalance := newBlock.balanceOf(prevBlock.RewardAddr)
		newBlock.Balances[prevBlock.RewardAddr] = winnerBalance
	}

	return &newBlock
}
