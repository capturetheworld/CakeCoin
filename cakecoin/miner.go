package main

import (
	"crypto/rsa"
	"fmt"

	"./utils"
	//"./emitter"
)

//Miner is a client but since there is no subclasses
type Miner struct {
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
	CurrentBlock                                             *Block
	MiningRounds                                             int
}

func newMiner(name string, Net FakeNet, startingBlock *Block, keyPair *rsa.PrivateKey) *Miner {
	var m Miner
	m.Net = Net
	m.Name = name

	if keyPair == nil {
		m.KeyPair = utils.GenerateKeypair()
	} else {
		m.KeyPair = keyPair
	}

	m.Address = utils.CalculateAddress(&m.KeyPair.PublicKey)
	m.Nonce = 0

	m.PendingOutgoingTransactions = make(map[string]*Transaction)
	m.PendingReceivedTransactions = make(map[string]*Transaction)
	m.Blocks = make(map[string]*Block)
	m.PendingBlocks = make(map[string]*Block)

	if startingBlock != nil {
		m.setGenesisBlock(startingBlock)
	}

	m.MiningRounds = NUM_ROUNDS_MINING

	return &m

}

func (m Miner) setGenesisBlock(startingBlock *Block) {
	if m.LastBlock != nil {
		fmt.Printf("Cannot set starting block for existing blockchain.")
	}

	m.LastConfirmedBlock = startingBlock
	m.LastBlock = startingBlock
	m.Blocks[string(startingBlock.id())] = startingBlock
}

/**
 * Starts listeners and begins mining.
 */
func (m Miner) initialize() {
	m.startNewSearch()

	this.on(Blockchain.START_MINING, this.findProof)
	this.on(Blockchain.POST_TRANSACTION, this.addTransaction)

	// setTimeout(() => this.emit(Blockchain.START_MINING), 0); //needs implementing
}

func (m Miner) startNewSearch() {
	m.currentBlock = Blockchain.makeBlock(m.address, m.lastBlock)

	// txSet.forEach((tx) => this.addTransaction(tx));

	// Start looking for a proof at 0.
	this.currentBlock.proof = 0
}


func (m Miner) announceProof() {
    m.Net.broadcast(PROOF_FOUND, m.CurrentBlock);
  }

func (m Miner) receiveBlock(s *Block) {
    m.Net.broadcast(PROOF_FOUND, m.CurrentBlock);
}


//no optional parameters
func (m Miner) findProof(oneAndDone bool) {
	if(oneAndDone == nil){
		oneAndDone = false
	}

	 pausePoint  := m.CurrentBlock.Proof + m.MiningRounds

    for (m.CurrentBlock.Proof < pausePoint) {
      if (m.currentBlock.hasValidProof() == true) {
		fmt.Printf("found proof for block %v",m.CurrentBlock.ChainLength)
		fmt.Printf(": %v\n",m.CurrentBlock.Proof)

        // this.log(`found proof for block ${this.currentBlock.chainLength}: ${this.currentBlock.proof}`);
        m.announceProof();
        m.receiveBlock(&m.CurrentBlock);
        m.startNewSearch();
        break;
	  }
	}

	m.CurrentBlock.Proof++;
	  
    // If we are testing, don't continue the search.
    if(oneAndDone == false) {
      // Check if anyone has found a block, and then return to mining.
      setTimeout(() => m.emit(Blockchain.START_MINING), 0);
    }
  }

