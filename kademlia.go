package kademlia

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"time"
)

type KademliaStorage interface {
	Put(key NodeID, value []byte) error
	Get(key NodeID) (value []byte, err error)
}

type Kademlia struct {
	routes    *RoutingTable
	Storage   KademliaStorage
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

// Generic RPC base
type RPCHeader struct {
	Sender    Contact
	NetworkID string
}

// Every RPC updates routing tables in Kademlia
func (k *Kademlia) HandleRPC(request RPCHeader, response *RPCHeader) error {
	if request.NetworkID != k.NetworkID {
		return errors.New(fmt.Sprintf("Expected Network ID %s, go %s", k.NetworkID, request.NetworkID))
	}

	// Update routing table for all incoming RPCs
	k.routes.Update(request.Sender)
	// Pong with sender
	response.Sender = k.routes.self

	return nil
}

func dialContact(contact Contact) (*rpc.Client, error) {
	connection, err := net.DialTimeout("tcp", contact.Address, 5*time.Second)
	if err != nil {
		return nil, err
	}

	return rpc.NewClient(connection), nil
}

func (k *Kademlia) Serve() error {
	rpc.Register(&KademliaCore{k})

	l, err := net.Listen("tcp", k.routes.self.Address)
	if err != nil {
		return err
	}

	go rpc.Accept(l)

	return nil
}

/*
 * KademliaCore
 * Handles RPC interactions between client/server
 */

type KademliaCore struct {
	kad *Kademlia
}
