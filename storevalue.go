package kademlia

type StoreValueRequest struct {
	RPCHeader
	Target NodeID
	Value []byte
}

func (k *Kademlia) NewStoreValueRequest(target NodeID, value []byte) StoreValueRequest {
	return StoreValueRequest{
		RPCHeader: RPCHeader{
			Sender:    k.routes.self,
			NetworkID: k.NetworkID,
		},
		Target: target,
		Value: value,
	}
}

type StoreValueResponse struct {
	RPCHeader
	Contacts Contacts
}

func (k *Kademlia) StoreValue(contact Contact, target NodeID, value []byte) ([]Contact, error) {
	client, err := dialContact(contact)
	if err != nil {
		return nil, err
	}

	req := k.NewStoreValueRequest(target, value)
	res := StoreValueResponse{}

	err = client.Call("KademliaCore.StoreValueRPC", &req, &res)
	if err != nil {
		return nil, err
	}

	return res.Contacts, nil
}

func (kc *KademliaCore) StoreValueRPC(req StoreValueRequest, res *StoreValueResponse) error {
	err := kc.kad.HandleRPC(req.RPCHeader, &res.RPCHeader)
	if err != nil {
		return err
	}

	p_err := kc.kad.Storage.Put(req.Target, req.Value[:])
	if p_err != nil {
		//log.Println(err)
		//panic("Read from values database failed")
		return p_err
	}

	res.Contacts = kc.kad.routes.FindClosest(req.Target, BucketSize)

	return nil
}
