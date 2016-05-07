package kademlia

type StoreValueRequest struct {
	CallHeader
	Target NodeID
	Value []byte
}

func (k *Kademlia) NewStoreValueRequest(target NodeID, value []byte) StoreValueRequest {
	return StoreValueRequest{
		CallHeader: CallHeader{
			Sender:    k.routes.self,
			NetworkID: k.NetworkID,
		},
		Target: target,
		Value: value,
	}
}

type StoreValueResponse struct {
	CallHeader
	Contacts Contacts
}

func (k *Kademlia) StoreValue(contact Contact, target NodeID, value []byte) ([]Contact, error) {
	client, err := k.Network.Connect(contact)
	if err != nil {
		return nil, err
	}

	req := k.NewStoreValueRequest(target, value)
	res := StoreValueResponse{}

	err = client.StoreValue(req, &res)
	if err != nil {
		return nil, err
	}

	return res.Contacts, nil
}

func (k *Kademlia) StoreValueHandler(req StoreValueRequest, res *StoreValueResponse) error {
	err := k.HandleCall(req.CallHeader, &res.CallHeader)
	if err != nil {
		return err
	}

	p_err := k.Storage.Put(req.Target, req.Value[:])
	if p_err != nil {
		//log.Println(err)
		//panic("Read from values database failed")
		return p_err
	}

	res.Contacts = k.routes.FindClosest(req.Target, BucketSize)

	return nil
}
