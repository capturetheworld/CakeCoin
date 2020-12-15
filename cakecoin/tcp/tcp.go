package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"../utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

// SimpleBlock is a simpler version of Block
type SimpleBlock struct {
	PrevBlockHash string
	id            int
	Timestampstr  string
	Transaction   int
	Hash          string
	ChainLength   int
}

//Blockchain
var bc []SimpleBlock

// SimpleBlockhainServer handles incoming concurrent SimpleBlocks
var SimpleBlockhainServer chan []SimpleBlock

//mutex lock
var mutex = &sync.Mutex{}

func main() {

	tcpPort := "8000" //specify the TCP Port

	err := godotenv.Load()
	if err != nil {
		log.Println("Loading error")
	}

	SimpleBlockhainServer = make(chan []SimpleBlock)

	MakeGenesis()

	// startup of TCP server
	server, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		log.Println("Can't listen on TCP Port")
	}
	log.Println("TCP Server Listening on port :", tcpPort)
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("Seems to have some type of server excepting error")
		}
		go tcpConnections(conn)
	}

}

func MakeGenesis() {
	t := time.Now()
	fmt.Println("\n\n\n*****************CREATING GENESIS BLOCK*****************")
	genesisSimpleBlock := NewBlock("", 0, t.String(), 0, "")

	spew.Dump(genesisSimpleBlock)
	fmt.Println("*******************CREATED GENESIS BLOCK***************\n\n\n")
	bc = append(bc, *genesisSimpleBlock)
}

//from TCP guides
func tcpConnections(conn net.Conn) {

	defer conn.Close()

	io.WriteString(conn, "Enter a new value:")

	scanner := bufio.NewScanner(conn)

	// take in Transaction value from Scanner
	go func() {
		for scanner.Scan() {
			Transaction, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}
			newSimpleBlock, err := AddNewSimpleBlock(bc[len(bc)-1], Transaction)
			if err != nil {
				log.Println("Can't create a new block for the blockchain")
				continue
			}
			if hasValidProof(newSimpleBlock, bc[len(bc)-1]) {
				newbc := append(bc, newSimpleBlock)
				replaceChain(newbc)
			}

			SimpleBlockhainServer <- bc
			io.WriteString(conn, "\nEnter a new Transaction:")
		}
	}()

	for _ = range SimpleBlockhainServer {
		spew.Dump(bc)
	}

}

// check chain length
func replaceChain(newSimpleBlocks []SimpleBlock) {
	mutex.Lock()
	if len(newSimpleBlocks) > len(bc) {
		bc = newSimpleBlocks
	}
	mutex.Unlock()
}

////////////////SIMPLE BLOCK METHODS//////////////////
func hashVal(SimpleBlock SimpleBlock) string {
	record := string(SimpleBlock.id) + SimpleBlock.Timestampstr + string(SimpleBlock.Transaction) + SimpleBlock.PrevBlockHash
	return hex.EncodeToString(utils.Hash(record))
}

// check for valid proof
func hasValidProof(newSB, oldSB SimpleBlock) bool {
	if oldSB.id+1 != newSB.id ||
		oldSB.Hash != newSB.PrevBlockHash ||
		hashVal(newSB) != newSB.Hash {
		return false
	}

	return true
}

func NewBlock(PBH string, ident int, ts string, tx int, hash string) *SimpleBlock {
	newBlock := SimpleBlock{PBH, ident, ts, tx, hash, 0}
	return &newBlock
}

func AddNewSimpleBlock(oldSimpleBlock SimpleBlock, Transaction int) (SimpleBlock, error) {

	var s SimpleBlock

	t := time.Now()

	s.id = oldSimpleBlock.id + 1
	s.Timestampstr = t.String()
	s.Transaction = Transaction
	s.PrevBlockHash = oldSimpleBlock.Hash
	s.Hash = hashVal(s)

	return s, nil
}
