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

type BlockChain struct {
	PowTarget      *big.Int
	CoinbaseAmount int
	DefaultTxFee   int
	ConfirmedDepth int
}

//Not sure if we need some of these functions as everything is visible anyway
/**
func GET_MISSING_BLOCK() string    { return MISSING_BLOCK }
func GET_POST_TRANSACTION() string { return POST_TRANSACTION }
func GET_PROOF_FOUND() string      { return PROOF_FOUND }
func GET_START_MINING() string     { return START_MINING }

func GET_NUM_ROUNDS_MINING() int { return NUM_ROUNDS_MINING }
// Configurable properties.
func GET_POW_TARGET()           { return Blockchain.cfg.powTarget }
func GET_COINBASE_AMT_ALLOWED() { return Blockchain.cfg.coinbaseAmount }
func GET_DEFAULT_TX_FEE()       { return Blockchain.cfg.defaultTxFee }
func GET_CONFIRMED_DEPTH()      { return Blockchain.cfg.confirmedDepth }
**/

//Ignore using function pointers for now

//implement defaults here instead?
func (bc *BlockChain) MakeGenesis(powLeadingZeroes int, coinbaseAmount int, defaultTxFee int, confirmedDepth int, clientBalanceMap map[*Client]int, startingBalances map[string]int) *Block {
	if clientBalanceMap != nil && startingBalances != nil {
		panic("Can only have either clientbalancemap or startingbalances")
	}
	//set config of blockchain?
	bc.CoinbaseAmount = coinbaseAmount
	bc.DefaultTxFee = defaultTxFee
	bc.ConfirmedDepth = confirmedDepth
	//set powtarget
	n := new(big.Int)
	n, err := n.SetString(POW_BASE_TARGET, 16)
	if !err {
		panic("can't set big int")
	}
	n.Rsh(n, uint(powLeadingZeroes))

	bc.PowTarget = n

	balances := make(map[string]int)
	if startingBalances != nil {
		balances = startingBalances
	}

	if clientBalanceMap != nil {
		for client, val := range clientBalanceMap {
			balances[client.Address] = val
		}
	}

	g := MakeBlock("", nil, bc.PowTarget, bc.CoinbaseAmount)

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

//probably accepts a json object
func deserializeBlock(o []byte) *Block {
	if !json.Valid(o) {
		panic("Input is not a valid json object for block")
	}

	var jsonBlock Block
	var b Block
	/**
	dec := json.NewDecoder(strings.NewReader(string(o)))
	if err := dec.Decode(&jsonBlock); err == io.EOF {
	} else if err != nil {
		log.Fatal(err)
	}
	**/
	err := json.Unmarshal(o, &jsonBlock)
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonBlock.Serialize())

	balances := make(map[string]int)
	chainLength := jsonBlock.ChainLength
	timestamp := jsonBlock.Timestamp

	if jsonBlock.IsGenesisBlock() {
		//fmt.Println("setting balances")
		//fmt.Println(jsonBlock.Balances)
		for client, amount := range jsonBlock.Balances {
			balances[client] = amount
		}
		b.Balances = balances
	} else {
		prevBlockHash := jsonBlock.PrevBlockHash
		proof := jsonBlock.Proof
		rewardAddr := jsonBlock.RewardAddr
		transactions := make(map[string]*Transaction)
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
