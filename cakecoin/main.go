package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	fakeNet := NewFakeNet()
	Alice := NewClient("Alice", fakeNet, nil, nil)
	Bob := NewClient("Bob", fakeNet, nil, nil)
	Mickey := newMiner("Mickey", fakeNet, nil, nil)
	Minnie := newMiner("Minnie", fakeNet, nil, nil)

	fakeNet.register(Alice, Bob, Mickey.MinerClient, Minnie.MinerClient)

	bc := BlockChain{}

	clientBalanceMap := make(map[*Client]int)
	clientBalanceMap[Alice] = 200
	clientBalanceMap[Bob] = 100
	clientBalanceMap[Mickey.MinerClient] = 50
	clientBalanceMap[Minnie.MinerClient] = 50

	g := bc.MakeGenesis(18, 25, 5, 6, clientBalanceMap, nil)
	fmt.Println(g.Serialize())
	outputs := []Output{Output{3.0, Bob.Address}, Output{4.0, Mickey.MinerClient.Address}}

	showBalances := func(c *Client) {
		fmt.Printf("Alice has %v gold.\n", c.LastBlock.balanceOf(Alice.Address))
		fmt.Printf("Bob has %v gold.\n", c.LastBlock.balanceOf(Bob.Address))
		fmt.Printf("Minnie has %v gold.\n", c.LastBlock.balanceOf(Minnie.MinerClient.Address))
		fmt.Printf("Mickey has %v gold.\n", c.LastBlock.balanceOf(Mickey.MinerClient.Address))
		fmt.Println()
	}

	timeout1 := func() {
		fmt.Printf("")
		fmt.Printf("Mickey has a chain of length %v\n", Mickey.CurrentBlock.ChainLength)

		fmt.Printf("")
		fmt.Printf("Minnie has a chain of length %v\n", Minnie.CurrentBlock.ChainLength)

		fmt.Printf("")
		fmt.Printf("Final balances (Minnie's perspective):\n")
		showBalances(Minnie.MinerClient)
		fmt.Printf("")
		fmt.Printf("Final balances (Alice's perspective):\n")
		showBalances(Alice)
		os.Exit(0)
	}
	DurationOfTime := time.Duration(10) * time.Second
	time.AfterFunc(DurationOfTime, timeout1)

	fmt.Printf("Initial balances:\n")
	showBalances(Alice)

	go Mickey.initialize()
	go Minnie.initialize()

	go Alice.postTransaction(outputs, 5)

	time.Sleep(20 * time.Second)
}
