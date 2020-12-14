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

	fakeNet.register(Alice, Bob, Mickey.MinerClient)

	bc := BlockChain{}

	clientBalanceMap := make(map[*Client]int)
	clientBalanceMap[Alice] = 200
	clientBalanceMap[Bob] = 100
	clientBalanceMap[Mickey.MinerClient] = 50

	g := bc.MakeGenesis(15, 25, 5, 6, clientBalanceMap, nil)
	fmt.Println(g.Serialize())
	fmt.Println(Alice.confirmedBalance())
	fmt.Println(Alice.availableGold())

	outputs := []Output{Output{3.0, Bob.Address}, Output{4.0, Mickey.MinerClient.Address}}

	Mickey.initialize()

	Alice.postTransaction(outputs, 5)

	time.AfterFunc(5*time.Second, timeout1)
}

func timeout1() {
	os.Exit(0)
}
