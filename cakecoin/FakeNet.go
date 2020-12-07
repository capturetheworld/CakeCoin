package main

import (
	"encoding/json"
	 "time"
	
)


type FakeNet struct {
	Clients map[string]*Client
	o2 interface{}
}


func (f FakeNet) register(clientList ...Client)  {
	for _, client := range clientList {
		f.Clients[client.Address] = &client
	}
}

func (f FakeNet) broadcast(msg string,o interface{}) {
	for address, _ := range f.Clients {
		f.sendMessage(address,msg,o)
	}
}


func (f FakeNet) recognizes(client Client) bool  {
	if val, ok := f.Clients[client.Address]; ok {
		return true
	}else{
		return false
	}
}

func (f FakeNet) sendMessage(addr string,msg string, o interface{}) {
	jsonByte, err := json.Marshal(o)
	f.o2 = json.Unmarshal(string(jsonByte))
	//needs setTimeout(() => client.emit(msg, o2), 0);
	var CLIENT = (f.Clients[addr]);
	time.AfterFunc(0, CLIENT.emit(msg,f.o2))

}




func NewFakeNet() *FakeNet{
	var f FakeNet
	f.Clients = make(map[string]*Client)

	return &f
}


