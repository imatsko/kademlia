package kademlia

type FindValueRequest struct {
	RPCHeader
	Target NodeID
}

func (k *Kademlia) NewFindValueRequest(target NodeID) FindValueRequest {
	return FindValueRequest{
		RPCHeader: RPCHeader{
			Sender:    k.routes.self,
			NetworkID: k.NetworkID,
		},
		Target: target,
	}
}

type FindValueResponse struct {
	RPCHeader
	Contacts Contacts
	Value    []byte
}

func (k *Kademlia) FindValue(contact Contact, target NodeID) ([]Contact, []byte, error) {
	client, err := dialContact(contact)
	if err != nil {
		return nil, nil, err
	}

	req := k.NewFindValueRequest(target)
	res := FindValueResponse{}

	err = client.Call("KademliaCore.FindValueRPC", &req, &res)
	if err != nil {
		return nil, nil, err
	}

	return res.Contacts, res.Value, nil
}

func (kc *KademliaCore) FindValueRPC(req FindValueRequest, res *FindValueResponse) error {
	err := kc.kad.HandleRPC(req.RPCHeader, &res.RPCHeader)
	if err != nil {
		return err
	}

	value, err := kc.kad.Storage.Get(req.Target)
	if err != nil {
		//log.Println(err)
		//panic("Read from values database failed")
		return err
	}

	if value != nil {
		res.Value = value
		return nil
	}

	res.Contacts = kc.kad.routes.FindClosest(req.Target, BucketSize)

	return nil
}
