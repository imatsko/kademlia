package kademlia

type PingRequest struct {
	CallHeader
}

func (k *Kademlia) NewPingRequest() PingRequest {
	return PingRequest{
		CallHeader{
			Sender:    k.routes.self,
			NetworkID: k.NetworkID,
		},
	}
}

type PingResponse struct {
	CallHeader
}

func (k *Kademlia) Ping(target Contact) error {
	client, err := k.Network.Connect(target)
	if err != nil {
		return err
	}

	req := k.NewPingRequest()
	res := PingResponse{}

	return client.Ping(req, &res)
}

func (k *Kademlia) PingHandler(req PingRequest, res *PingResponse) error {
	return k.HandleCall(req.CallHeader, &res.CallHeader)
}
