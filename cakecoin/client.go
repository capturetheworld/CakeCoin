package main

import (
	"crypto/rsa"
	"fmt"

	//"./utils"
	//"./emitter"
	"github.com/Stan/168proj/cakecoin/utils"
	"github.com/chuckpreslar/emission"
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
	Net                                                      *FakeNet `json:"-"`
	KeyPair                                                  *rsa.PrivateKey
	PendingOutgoingTransactions, PendingReceivedTransactions map[string]*Transaction
	Blocks                                                   map[string]*Block
	PendingBlocks                                            map[string][]*Block
	LastBlock                                                *Block
	LastConfirmedBlock                                       *Block
	ReceiveBlock                                             *Block
	BlockChain                                               *BlockChain
	Emitter                                                  *emission.Emitter `json:"-"`
}

func (c *Client) setGenesisBlock(startingBlock *Block) {
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
	for _, tx := range c.PendingOutgoingTransactions {
		pendingSpent = pendingSpent + tx.TotalOutput()
	}

	return int(c.confirmedBalance()) - pendingSpent

}

func (c *Client) postTransaction(outputs []Output, fee int) *Transaction {
	total := fee
	for _, output := range outputs {
		total += output.Amount
	}
	if total > c.availableGold() {
		panic(`Account doesn't have enough gold for transaction`)
	}

	tx := NewTransaction(c.Address, c.Nonce, &c.KeyPair.PublicKey, nil, fee, outputs)

	tx.Sign(c.KeyPair)
	fmt.Printf("Alice signs and it has a %v sig\n", tx.ValidSignature())

	c.PendingOutgoingTransactions[string(tx.Id())] = tx

	c.Nonce++

	c.Net.broadcast(POST_TRANSACTION, tx)

	return tx
}

func (c *Client) receiveBlock(b *Block, bstr string) *Block {
	//fmt.Printf("%s receiving Block\n", c.Name)
	block := b

	if b == nil {
		fmt.Println("Deseralize")
		block = c.BlockChain.deserializeBlock([]byte(bstr))
	}
	if _, ok := c.Blocks[string(block.id())]; ok {
		return nil
	}
	if !block.hasValidProof() && !block.IsGenesisBlock() {
		fmt.Printf("Block %v does not have a valid proof\n", string(block.id()))
		return nil
	}

	//make sure that the if statement after this actually sets it
	var prevBlock *Block = nil

	prevBlock, ok := c.Blocks[string(block.PrevBlockHash)]
	if !ok {
		if prevBlock == nil || !prevBlock.IsGenesisBlock() {
			stuckBlocks, ok := c.PendingBlocks[string(block.PrevBlockHash)]
			if !ok {
				c.requestMissingBlock(*block)
				//magic number here
				stuckBlocks = make([]*Block, 10)
			}
			stuckBlocks = append(stuckBlocks, block)
			c.PendingBlocks[string(block.PrevBlockHash)] = stuckBlocks
			return nil
		}
	}

	if !block.IsGenesisBlock() {
		if !block.rerun(prevBlock) {
			return nil
		}
	}
	//may be the cycle right here
	c.Blocks[string(block.id())] = block

	if c.LastBlock.ChainLength < block.ChainLength {
		c.LastBlock = block
		c.setLastConfirmed()
	}
	//magic number
	unstuckBlocks := make([]*Block, 0)
	if val, ok := c.PendingBlocks[string(block.id())]; ok {
		unstuckBlocks = val
	}
	delete(c.PendingBlocks, string(block.id()))

	for _, uBlock := range unstuckBlocks {
		fmt.Printf("processing unstuck block %v", string(block.id()))
		c.receiveBlock(uBlock, "")
	}

	return block

}

func (c Client) requestMissingBlock(b Block) {
	fmt.Printf("Asking for missing block: %v", b.PrevBlockHash)
	var msg = Message{c.Address, b.PrevBlockHash}
	c.Net.broadcast(MISSING_BLOCK, msg)
}

func (c Client) resendPendingTransactions() {
	for _, tx := range c.PendingOutgoingTransactions {
		//EMIT THE TRANSACTION
		c.Net.broadcast(POST_TRANSACTION, tx)
	}
}

func (c Client) provideMissingBlock(msg Message) {
	//asking for missing message, just use prevblockhashfornow
	if val, ok := c.Blocks[string(msg.PrevBlockHash)]; ok {
		fmt.Printf("Providing missing block %v", string(msg.PrevBlockHash))
		block := val
		//EMIT MESSAGE WITH BLOCk
		c.Net.sendMessage(msg.Address, PROOF_FOUND, block)
	}
}

//func setLastConfirmed
func (c *Client) setLastConfirmed() {
	block := c.LastBlock
	confirmedBlockHeight := block.ChainLength - CONFIRMED_DEPTH
	if confirmedBlockHeight < 0 {
		confirmedBlockHeight = 0
	}
	for block.ChainLength > confirmedBlockHeight {
		block = c.Blocks[string(block.PrevBlockHash)]
	}
	c.LastConfirmedBlock = block

	for id, tx := range c.PendingOutgoingTransactions {
		if c.LastConfirmedBlock.contains(tx) {
			delete(c.PendingOutgoingTransactions, id)
		}
	}
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
	name := c.Address[0:10]
	if len(c.Name) > 0 {
		name = c.Name
	}

	fmt.Printf("    %v", name)
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

func (c *Client) receive(b *Block) {
	c.receiveBlock(b, "")
}

func NewClient(name string, Net *FakeNet, startingBlock *Block, keyPair *rsa.PrivateKey) *Client {
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
	c.PendingBlocks = make(map[string][]*Block)

	if startingBlock != nil {
		c.setGenesisBlock(startingBlock)
	}

	c.Emitter = emission.NewEmitter()
	c.Emitter.On(PROOF_FOUND, c.receive)
	c.Emitter.On(MISSING_BLOCK, c.provideMissingBlock)
	return &c
}
