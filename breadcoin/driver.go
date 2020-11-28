package main


import (
	"fmt"
	// Blockchain "./blockchain"
	// Block "./block"
	// Client "./client"
	// Miner "./miner"
	// Transaction "./transaction"
	// FakeNet "./fakeNet"
)

var fakeNet = FakeNet{}

//Clients
var alice = Client{name: "Alice", net: fakeNet}
var bob = Client{name: "Bob", net: fakeNet}
var charlie = Client{name: "Charlie", net: fakeNet}

//Miners
var minnie = Miner{name: "Minnie", net: fakeNet}
var mickey = Miner{name: "Mickey", net: fakeNet}


var donald =  Miner{name: "Donald", net: fakeNet, startingBlock: genesis, miningRounds: 3000};

func main() {
	fmt.Printf("Starting simulation.  This may take a moment...")

	fmt.Printf("Initial balances:")
	showBalances(alice)
	
	fakeNet.register(alice, bob, charlie, minnie, mickey)

	minnie.initialize()
	mickey.initialize()

	

}