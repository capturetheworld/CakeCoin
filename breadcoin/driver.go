package main

import (
	"fmt"
	"os"
	"time"
)

var fakeNet = FakeNet{}

//Clients
var alice = Client{name: "Alice", net: fakeNet}
var bob = Client{name: "Bob", net: fakeNet}
var charlie = Client{name: "Charlie", net: fakeNet}

//Miners
var minnie = Miner{name: "Minnie", net: fakeNet}
var mickey = Miner{name: "Mickey", net: fakeNet}

// var genesis  //to implement

var donald = Miner{name: "Donald", net: fakeNet, startingBlock: genesis, miningRounds: 3000}

func notmain() {
	fmt.Printf("Starting simulation.  This may take a moment...")

	fmt.Printf("Initial balances:")
	showBalances(alice)

	fakeNet.register(alice, bob, charlie, minnie, mickey)

	minnie.initialize()
	mickey.initialize()

	fmt.Printf("Alice is transferring 40 gold to %v", bob.address)
	alice.postTransaction(40, bob.address)

	time.AfterFunc(2*time.Second, timeout1())

	time.AfterFunc(5*time.Second, timeout2())

}

func timeout1() {
	fmt.Printf("")
	fmt.Printf("Starting simulation.  This may take a moment...")
	fakeNet.register(donald)
	donald.initialize()

}

func timeout2() {

	fmt.Printf("")
	fmt.Printf("Minnie has a chain of length %v", minnie.currentBlock.chainLength)

	fmt.Printf("")
	fmt.Printf("Mickey has a chain of length %v", mickey.currentBlock.chainLength)

	fmt.Printf("")
	fmt.Printf("Donald has a chain of length %v", donald.currentBlock.chainLength)

	fmt.Printf("")
	fmt.Printf("Final balances (Minnie's perspective):")
	showBalances(minnie)

	fmt.Printf("")
	fmt.Printf("Final balances (Alice's perspective):")
	showBalances(alice)

	fmt.Printf("")
	fmt.Printf("Final balances (Donald's perspective):")
	showBalances(donald)

	os.Exit(0)
}
