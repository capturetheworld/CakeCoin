package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"strings"
)

//Include constants here
const MISSING_BLOCK = "MISSING_BLOCK"
const POST_TRANSACTION = "POST_TRANSACTION"
const PROOF_FOUND = "PROOF_FOUND"
const START_MINING = "START_MINING"

const NUM_ROUNDS_MINING = 2000

const POW_BASE_TARGET = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
const POW_LEADING_ZEROES = 15

const COINBASE_AMT_ALLOWED = 25
const DEFAULT_TX_FEE = 1

const CONFIRMED_DEPTH = 6

//Ignore using function pointers for now
/**
func MakeGenesis(powLeadingZeroes int, coinbaseAmount float64, defaultTxFee float64, confirmedDepth int, clientBalanceMap map[Client]float64, startingBalances map[string]float64) *Block {
	if clientBalanceMap != nil && startingBalances != nil {
		panic("Can only have either clientbalancemap or startingbalances")
	}
	//set config of blockchain?
	//set powtarget

	balances := make(map[string]float64)
	if startingBalances != nil {
		balances = startingBalances
	}

	if clientBalanceMap != nil {
		for client, val := range clientBalanceMap {
			balances[client.Address] = val
		}
	}

	g := MakeBlock()
	for address, val := range balances {
		g.Balances[address] = val
	}

	if clientBalanceMap != nil {
		for client, _ := range clientBalanceMap {
			client.setGenesisBlock(g)
		}
	}

	return g
}
**/
//probably accepts a json object
func deserializeBlock(o []byte) *Block {
	if !json.Valid(o) {
		panic("Input is not a valid json object for block")
	}

	var jsonBlock Block
	var b Block
	dec := json.NewDecoder(strings.NewReader(string(o)))
	if err := dec.Decode(&jsonBlock); err == io.EOF {
	} else if err != nil {
		log.Fatal(err)
	}

	fmt.Println(jsonBlock.Serialize())

	balances := make(map[string]float64)
	chainLength := jsonBlock.ChainLength
	timestamp := jsonBlock.Timestamp

	if jsonBlock.IsGenesisBlock() {
		for client, amount := range jsonBlock.Balances {
			balances[client] = amount
		}
	} else {
		prevBlockHash := jsonBlock.PrevBlockHash
		proof := jsonBlock.Proof
		rewardAddr := jsonBlock.RewardAddr
		transactions := make(map[string]Transaction)
		if jsonBlock.Transactions != nil {
			for id, tx := range jsonBlock.Transactions {
				transactions[id] = tx
			}
		}
		//GOTTA FIX THIS WHEN YOU IMPLEMENT CONSTANTS
		b = *NewBlock(rewardAddr, nil, nil, 25)
		b.ChainLength = chainLength
		b.Timestamp = timestamp
		b.PrevBlockHash = prevBlockHash
		b.Proof = proof
		b.Balances = balances
		b.Transactions = transactions
	}
	return &b
}

func MakeBlock(s string, b *Block, i *big.Int, c int) *Block {
	return NewBlock(s, b, i, c)
}

func MakeTransaction(o []byte) *Transaction {
	if !json.Valid(o) {
		panic("Input is not a valid json object for transaction")
	}

	var jsonTransaction Transaction
	dec := json.NewDecoder(strings.NewReader(string(o)))
	if err := dec.Decode(&jsonTransaction); err == io.EOF {
	} else if err != nil {
		log.Fatal(err)
	}
	return NewTransaction(jsonTransaction.From, jsonTransaction.Nonce, jsonTransaction.PubKey, jsonTransaction.Sig, jsonTransaction.Fee, jsonTransaction.Outputs)
}
