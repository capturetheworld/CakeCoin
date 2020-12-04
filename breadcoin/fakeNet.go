package main

import {
	"encoding/json"
}



type FakeNet struct {
	Clients map[[]byte]Client
}


func (f FakeNet) register(clientList ...Client)  {
	for _, client := range clientList {
		f.Clients[client.Address] = client
	}
}

func (f FakeNet) broadcast(msg string,o interface{}) {
	for address, client := range f.Clients {
		f.sendMessage(address,msg,o)
	}
}


func (f FakeNet) recognizes(client Client) bool  {
	if val, ok := clientList[client.Address]; ok {
		return true
	}
	else{
		return false
	}
}

func (f FakeNet) sendMessage(addr string,msg string, o interface{}) {
	jsonByte, err := json.Marshal(o)
	f.o2 = json.Unmarshal(string(jsonByte))
	//needs setTimeout(() => client.emit(msg, o2), 0);
}


func NewFakeNet(){
	var f FakeNet
	f.Clients = make(map[[]byte]Client)

	return &f
}




(