package main

import (
	"crypto/rsa"
	"fmt"

	"./utils"
	//"./emitter"
)

//Message is
type Message struct {
	Address       string
	PrevBlockHash []byte
}

//Client is
type Client struct {
	Name, Address                                            string
	Nonce                                                    int
	Net                                                      FakeNet
	KeyPair                                                  *rsa.PrivateKey
	PendingOutgoingTransactions, PendingReceivedTransactions map[string]*Transaction
	Blocks                                                   map[string]*Block
	PendingBlocks                                            map[string]*Block
	LastBlock                                                *Block
	LastConfirmedBlock                                       *Block
	ReceiveBlock                                             *Block
}

func (c Client) setGenesisBlock(startingBlock *Block) {
	if c.LastBlock != nil {
		fmt.Printf("Cannot set starting block for existing blockchain.")
	}

	c.LastConfirmedBlock = startingBlock
	c.LastBlock = startingBlock
	c.Blocks[string(startingBlock.id())] = startingBlock
}

func (c Client) confirmedBalance() int {
	return c.LastConfirmedBlock.balanceOf(c.Address)
}

func (c Client) availableGold() int {
	var pendingSpent int = 0
	for id, tx := range c.PendingOutgoingTransactions {
		pendingSpent = pendingSpent + tx.TotalOutput()
	}

	return int(c.confirmedBalance()) - pendingSpent

}

func (c Client) postTransaction() *Transaction {
	return nil
}

func (c Client) receiveBlock(b *Block) *Block { //needs finishing, figure out how to return null
	b = deserializeBlock(b)

	if val, ok := c.Blocks[b.ID]; ok {
		return nil
	}
	return nil

}

func (c Client) requestMissingBlock(b Block) {
	fmt.Printf("Asking for missing block: %v", b.PrevBlockHash)
	var msg = Message{c.Address, b.PrevBlockHash}
	c.Net.broadcast(MISSING_BLOCK, msg)
}

func (c Client) showAllBalances(b Block) {
	fmt.Printf("Showing balances:")

	for id, balance := range c.LastConfirmedBlock.Balances {
		fmt.Printf("    %v", id)
		fmt.Printf("    %v", balance)
		fmt.Println("")
	}
}

func (c Client) log(msg Message) {
	if len(c.Name) > 0 {
		name := c.Name
	} else {
		name := c.Address[0:10]
	}

	fmt.Printf("    %v", c.Name)
	fmt.Printf("    %v", msg)

}

func (c Client) showBlockchain() {
	block := c.LastBlock
	fmt.Println("BLOCKCHAIN:")
	for block != nil {
		fmt.Println(block.id())
		block = c.Blocks[string(block.PrevBlockHash)]
	}
}

func NewClient(name string, Net FakeNet, startingBlock *Block, keyPair *rsa.PrivateKey) *Client {
	var c Client
	c.Net = Net
	c.Name = name

	if keyPair == nil {
		c.KeyPair = utils.GenerateKeypair()
	} else {
		c.KeyPair = keyPair
	}

	c.Address = utils.CalculateAddress(&c.KeyPair.PublicKey)
	c.Nonce = 0

	c.PendingOutgoingTransactions = make(map[string]*Transaction)
	c.PendingReceivedTransactions = make(map[string]*Transaction)
	c.Blocks = make(map[string]*Block)
	c.PendingBlocks = make(map[string]*Block)

	if startingBlock != nil {
		c.setGenesisBlock(startingBlock)
	}

	return &c

}
