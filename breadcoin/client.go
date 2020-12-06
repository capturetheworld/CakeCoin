package main

import (
	"crypto/rsa"
	"fmt"
	"./utils"
	//"./emitter"
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
	LastBlock                                                *Block
	LastConfirmedBlock                                       *Block
	ReceiveBlock                                             *Block
}

func (c Client) setGenesisBlock(startingBlock *Block) {
	if (c.LastBlock != nil) {
		fmt.Printf("Cannot set starting block for existing blockchain.")
	}

	c.LastConfirmedBlock = startingBlock
	c.LastBlock = startingBlock
	c.Blocks[startingBlock.id()] = startingBlock
}

func (c Client) confirmedBalance() float64 {
	return c.LastConfirmedBlock.balanceOf(c.Address)
}

func (c Client) availableGold() float64 {
	var  pendingSpent float64 = 0.0
	for id, tx := range c.PendingOutgoingTransactions {
		pendingSpent += tx.TotalOutput()
	}

	return c.confirmedBalance() - pendingSpent

}

func (c Client) postTransaction() *Transaction {  //to implement, contains default parameter
	return nil
}

func (c Client) receiveBlock(b Block) *Block { //needs finishing, figure out how to return null
	b = Blockchain.deserializeBlock(b)

	if val, ok := c.Blocks[b.ID]; ok {
		return nil
	}
	return nil

}

func (c Client) requestMissingBlock(b Block) {
	fmt.Printf("Asking for missing block: %v", b.PrevBlockHash)
	var msg = Message{c.Address, b.PrevBlockHash}
	c.Net.broadcast(Blockchain.MISSING_BLOCK(), msg)
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
	if (len(c.Name) > 0) { 
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
		block = c.Blocks[block].PrevBlockHash //unsure need to fix
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
