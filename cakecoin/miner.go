package main

import (
	"crypto/rsa"
	"fmt"
	//"./utils"
	//"./emitter"
)

//Miner is a client but since there is no subclasses
type Miner struct {
	MinerClient  *Client
	MiningRounds int
	CurrentBlock *Block
}

func newMiner(name string, Net *FakeNet, startingBlock *Block, keyPair *rsa.PrivateKey) *Miner {
	c := NewClient(name, Net, startingBlock, keyPair)
	var m Miner
	m.MinerClient = c
	m.MiningRounds = NUM_ROUNDS_MINING

	return &m
}

func (m *Miner) setGenesisBlock(startingBlock *Block) {
	m.MinerClient.setGenesisBlock(startingBlock)
}

func (m *Miner) initialize() {
	m.startNewSearch(nil)

	addTx := func(trans *Transaction) {
		m.addTransaction(trans, "")
	}

	m.MinerClient.Emitter.On(START_MINING, m.findProof)
	m.MinerClient.Emitter.On(POST_TRANSACTION, addTx)

	m.MinerClient.Emitter.Emit(START_MINING)
}

func (m *Miner) startNewSearch(txSet map[*Transaction]int) {
	fmt.Println("Starting new search")
	m.CurrentBlock = m.MinerClient.BlockChain.MakeBlock(m.MinerClient.Address, m.MinerClient.LastBlock, nil, nil)

	transSet := make(map[*Transaction]int)
	if txSet != nil {
		transSet = txSet
	}
	// txSet.forEach((tx) => this.addTransaction(tx));
	for key, _ := range transSet {
		m.addTransaction(key, "")
	}
	// Start looking for a proof at 0.
	m.CurrentBlock.Proof = 0
}

//no optional parameters
func (m *Miner) findProof() {
	pausePoint := m.CurrentBlock.Proof + m.MiningRounds

	for m.CurrentBlock.Proof < pausePoint {
		if m.CurrentBlock.hasValidProof() {
			fmt.Printf("found proof for block %v", m.CurrentBlock.ChainLength)
			fmt.Printf(": %v\n", m.CurrentBlock.Proof)

			// this.log(`found proof for block ${this.currentBlock.chainLength}: ${this.currentBlock.proof}`);
			m.announceProof()
			m.receiveBlock(m.CurrentBlock, "")
			m.startNewSearch(nil)
			break
		}
		m.CurrentBlock.Proof++
	}

	m.MinerClient.Emitter.Emit(START_MINING)

}

func (m *Miner) announceProof() {
	m.MinerClient.Net.broadcast(PROOF_FOUND, m.CurrentBlock)
}

func (m *Miner) receiveBlock(b *Block, blockStr string) {
	block := m.MinerClient.receiveBlock(b, blockStr)
	if block == nil {
		return
	}

	if m.CurrentBlock != nil && block.ChainLength >= m.CurrentBlock.ChainLength {
		fmt.Println("cutting over to new chain")
		txSet := m.syncTransactions(block)
		m.startNewSearch(txSet)
	}
}

func (m *Miner) syncTransactions(nb *Block) map[*Transaction]int {
	cb := m.CurrentBlock
	cbTxs := make(map[*Transaction]int)
	nbTxs := make(map[*Transaction]int)

	for nb.ChainLength > cb.ChainLength {
		for _, tx := range nb.Transactions {
			nbTxs[tx] = 0
			nb = m.MinerClient.Blocks[string(nb.PrevBlockHash)]
		}
	}

	for cb != nil && string(cb.id()) != string(nb.id()) {
		for _, tx := range cb.Transactions {
			cbTxs[tx] = 0
		}
		for _, tx := range nb.Transactions {
			nbTxs[tx] = 0
		}
		cb = m.MinerClient.Blocks[string(cb.PrevBlockHash)]
		nb = m.MinerClient.Blocks[string(nb.PrevBlockHash)]
	}

	for _, tx := range nb.Transactions {
		delete(cbTxs, tx)
	}

	return cbTxs
}

func (m *Miner) addTransaction(tx *Transaction, txStr string) bool {
	//may cause issues
	var trans *Transaction = tx
	if tx == nil {
		trans = MakeTransaction([]byte(txStr))
	}
	return m.CurrentBlock.addTransaction(trans, 0)
}

func (m *Miner) postTransaction(outputs []Output, fee int) bool {
	tx := m.MinerClient.postTransaction(outputs, fee)
	return m.addTransaction(tx, "")
}
