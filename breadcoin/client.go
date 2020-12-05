package main

import (
	"crypto/rsa"
	"fmt"
)

type Message struct {
	Address       string
	PrevBlockHash []byte
}

type Client struct {
	Name, Address                                            string
	Nonce                                                    int
	Net                                                      FakeNet
	KeyPair                                                  *rsa.PrivateKey
	PendingOutgoingTransactions, PendingReceivedTransactions map[string]*Transaction
	Blocks                                                   map[string]*Block
	PendingBlocks                                            map[string]*Block
	LastBlock                                                Block
	LastConfirmedBlock                                       Block
	ReceiveBlock                                             Block
}

func (c Client) setGenesisBlock(startingBlock *Block) {
	if c.LastBlock {
		fmt.Printf("Cannot set starting block for existing blockchain.")
	}

	c.LastConfirmedBlock = startingBlock
	c.LastBlock = startingBlock
	c.Blocks[startingBlock.id()] = startingBlock
}

// func (c Client) confirmedBalance() float64 {
// 	return c.LastConfirmedBlock.balanceOf(c.Address)
// }

// func (c Client) availableGold() float64 {
// 	pendingSpent := 0
// 	for id, tx := range c.pendingOutgoingTransactions {
// 		pendingSpent += tx.TotalOutput
// 	}

// 	return c.confirmedBalance() - pendingSpent

// }

// func (c Client) postTransaction() Transaction {} //to implement, contains default parameter

// func (c Client) receiveBlock(b Block) Block { //needs finishing, figure out how to return null
// 	b = Blockchain.deserializeBlock(b)

// 	if val, ok := c.blocks[b.ID]; ok {
// 		return null
// 	}

// }

// func (c Client) requestMissingBlock(b Block) {
// 	fmt.Printf("Asking for missing block: %v", b.PrevBlockHash)
// 	var msg = Message{c.Address, b.prevBlockHash}
// 	c.net.broadcast(Blockchain.MISSING_BLOCK, msg)
// }

// func (c Client) showAllBalances(b Block) {
// 	fmt.Printf("Showing balances:")

// 	for id, balance := range c.LastConfirmedBlock.balances {
// 		fmt.Printf("    %v", id)
// 		fmt.Printf("    %v", balance)
// 		fmt.Println("")
// 	}
// }

// func (c Client) log(msg Message) {

// 	name := this.name || this.address.substring(0, 10)
// 	fmt.Printf("    %v", name)
// 	fmt.Printf("    %v", Message)

// }

// func (c Client) showBlockchain() {
// 	block := c.LastBlock
// 	fmt.Println("BLOCKCHAIN:")
// 	for block != Nil {
// 		fmt.Println(block.ID)
// 		block = c.Blocks.block.PrevBlockHash
// 	}
// }

// func NewClient(name string, Net FakeNet, startingBlock *Block, keyPair *rsa.PrivateKey) *Client {
// 	var c Client
// 	c.Net = FakeNet
// 	c.Name = name

// 	if keyPair == Nil {
// 		c.KeyPair = utils.GenerateKeyPair()
// 	} else {
// 		c.KeyPair = keyPair
// 	}

// 	c.Address = utils.CalculateAddress(c.keyPair.public)
// 	c.Nonce = 0

// 	c.PendingOutgoingTransactions = make(map[[]byte]Transaction)
// 	c.PendingReceivedTransactions = make(map[[]byte]Transaction)
// 	c.Blocks = make(map[[]byte]Block)
// 	PendingBlocks = make(map[[]byte]Block)

// 	if startingBlock {
// 		c.setGenesisBlock(startingBlock)
// 	}

// 	return &c

// }
