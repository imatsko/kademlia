package kademlia

import (
	"errors"
	"fmt"
)

type KademliaStorage interface {
	Put(key NodeID, value []byte) error
	Get(key NodeID) (value []byte, err error)
}

type KademliaNetwork interface {
	Connect(contact Contact) (client KademliaClient, err error)
}

type KademliaClient interface {
	FindNode(req FindNodeRequest, res *FindNodeResponse) error
	FindValue(req FindValueRequest, res *FindValueResponse) error
	Ping(req PingRequest, res *PingResponse) error
	StoreValue(req StoreValueRequest, res *StoreValueResponse) error
}

type KademliaNodeHandler interface {
	FindNodeHandler(req FindNodeRequest, res *FindNodeResponse) error
	FindValueHandler(req FindValueRequest, res *FindValueResponse) error
	PingHandler(req PingRequest, res *PingResponse) error
	StoreValueHandler(req StoreValueRequest, res *StoreValueResponse) error
}

type Kademlia struct {
	routes    *RoutingTable
	Storage   KademliaStorage
	Network   KademliaNetwork
	NetworkID string
}

func NewKademlia(self Contact, networkID string) *Kademlia {
	ret := &Kademlia{
		routes:    NewRoutingTable(self),
		Storage:  nil,
		NetworkID: networkID,
	}
	return ret
}

// Generic call base
type CallHeader struct {
	Sender    Contact
	NetworkID string
}

// Every call updates routing tables in Kademlia
func (k *Kademlia) HandleCall(request CallHeader, response *CallHeader) error {
	if request.NetworkID != k.NetworkID {
		return errors.New(fmt.Sprintf("Expected Network ID %s, go %s", k.NetworkID, request.NetworkID))
	}

	// Update routing table for all incoming RPCs
	k.routes.Update(request.Sender)
	// Pong with sender
	response.Sender = k.routes.self

	return nil
}
