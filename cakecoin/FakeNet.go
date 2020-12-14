package main

type FakeNet struct {
	Clients map[string]*Client
	o2      interface{}
}

func (f *FakeNet) register(clientList ...*Client) {
	for _, client := range clientList {
		f.Clients[client.Address] = client
	}
}

func (f *FakeNet) broadcast(msg string, o interface{}) {
	for address, _ := range f.Clients {
		//fmt.Println(address)
		//fmt.Println(f.Clients[address])
		f.sendMessage(address, msg, o)
	}
}

func (f *FakeNet) recognizes(client Client) bool {
	if _, ok := f.Clients[client.Address]; ok {
		return true
	} else {
		return false
	}
}

func (f *FakeNet) sendMessage(addr string, msg string, o interface{}) {
	/**
	jsonByte, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	o2 := o
	err = json.Unmarshal(jsonByte, &o2)
	if err != nil {
		panic(err)
	}
	**/
	client := f.Clients[addr]
	//fmt.Printf("Sending %s to %s\n", msg, client.Name)
	client.Emitter.Emit(msg, o)
}

func NewFakeNet() *FakeNet {
	var f FakeNet
	f.Clients = make(map[string]*Client)

	return &f
}
