package main

//Ignore using function pointers for now
func MakeGenesis(powLeadingZeroes int, coinbaseAmount float64, defaultTxFee float64, confirmedDepth int, clientBalanceMap map[Client]float64, startingBalances map[string]float64) *Block {
	if clientBalanceMap != nil && startingBalances != nil {
		panic("Can only have either clientbalancemap or startingbalances")
	}
	//set config of blockchain?
	//set powtarget

	balances := make(map[string]float64)
	if startingBalances != nil {
		balances = startingBalances
	}

	if clientBalanceMap != nil {
		for client, val := range clientBalanceMap {
			balances[Client.Address] = val
		}
	}

	g = MakeBlock()
	for address, val := range balances {
		g.Balances[address] = val
	}

	if clientBalanceMap != nil {
		for client, _ := range clientBalanceMap {
			client.setGenesisBlock(g)
		}
	}

	return g
}

func MakeBlock() *Block {
	return Block
}
