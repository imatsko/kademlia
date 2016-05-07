package kademlia

import (
	"net/rpc"
	"time"
	"net"
)

type RPCNode struct {
	kad KademliaNodeHandler
}

type RPCNodeCore struct {
	kad KademliaNodeHandler
}

func NewRPCNode(handler KademliaNodeHandler) *RPCNode {
	return &RPCNode{handler}
}

func (n *RPCNode) Serve(address string) error {
	rpc.Register(&RPCNodeCore{n.kad})

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	go rpc.Accept(l)

	return nil
}

func (n *RPCNodeCore) FindValueRPC(req FindValueRequest, res *FindValueResponse) error {
	return n.kad.FindValueHandler(req, res)
}

func (n *RPCNodeCore) FindNodeRPC(req FindNodeRequest, res *FindNodeResponse) error {
	return n.kad.FindNodeHandler(req, res)
}

func (n *RPCNodeCore) PingRPC(req PingRequest, res *PingResponse) error {
	return n.kad.PingHandler(req, res)
}

func (n *RPCNodeCore) StoreValueRPC(req StoreValueRequest, res *StoreValueResponse) error {
	return n.kad.StoreValueHandler(req, res)
}

type RPCNetwork struct {}

func NewRPCNetwork() *RPCNetwork {
	return new(RPCNetwork)
}

func (n *RPCNetwork) Connect(contact Contact) (c KademliaClient, err error) {
	RPCclient := new(RPCClientConnection)
	connection, err := net.DialTimeout("tcp", contact.Address, 5*time.Second)
	if err != nil {
		return nil, err
	}

	RPCclient.client = rpc.NewClient(connection)
	return RPCclient, nil
}

type RPCClientConnection struct {
	client *rpc.Client
}

func (c *RPCClientConnection) FindNode(req FindNodeRequest, res *FindNodeResponse) (err error) {
	err = c.client.Call("RPCNodeCore.FindNodeRPC", req, res)
	return
}

func (c *RPCClientConnection) FindValue(req FindValueRequest, res *FindValueResponse) (err error) {
	err = c.client.Call("RPCNodeCore.FindValueRPC", req, res)
	return
}

func (c *RPCClientConnection) Ping(req PingRequest, res *PingResponse) (err error) {
	err = c.client.Call("RPCNodeCore.PingRPC", req, res)
	return
}

func (c *RPCClientConnection) StoreValue(req StoreValueRequest, res *StoreValueResponse) (err error) {
	err = c.client.Call("RPCNodeCore.StoreValueRPC", req, res)
	return
}
