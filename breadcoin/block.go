package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/user/168proj/breadcoin/utils"
)

type Transaction struct {
	temp string
	id   string
}

type Block struct {
	PrevBlockHash  []byte
	Target         *big.Int
	Balances       map[string]float64
	NextNonce      map[string]int
	Transactions   map[string]Transaction
	ChainLength    int
	Timestamp      time.Time
	RewardAddr     string
	CoinbaseReward int
}

func (b Block) hashVal() []byte {
	return utils.Hash(b.serialize())
}

func (b Block) totalRewards() int {
	return b.CoinbaseReward
}

func (b Block) balanceOf(addr string) float64 {
	return b.Balances[addr]
}

func (b Block) isGenesisBlock() bool {
	return b.ChainLength == 0
}

func (b Block) hasValidProof() bool {
	h := b.hashVal()
	n := new(big.Int)
	n, err := n.SetString(string(h), 16)
	if !err {
		fmt.Println("Could not set hash of block to big Int")
		return false
	}
	return n.Cmp(b.Target) < 0
}

func (b Block) serialize() string {
	jsonByte, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return string(jsonByte)
}

func (b Block) id() []byte {
	return b.hashVal()
}

//UNOKENTNAEOITG! IMPLEMENT!!!!!!!!!!!!!!!!!!
func addTransaction(tx Transaction, client int) bool {
	return true
}

func New(rewardAddr string, prevBlock *Block, target *big.Int, coinbaseReward int) *Block {
	var prevBlockHash []byte = nil
	balances := make(map[string]float64)
	nextNonce := make(map[string]int)
	transactions := make(map[string]Transaction)
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
