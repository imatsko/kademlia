package kademlia

type FindValueRequest struct {
	CallHeader
	Target NodeID
}

func (k *Kademlia) NewFindValueRequest(target NodeID) FindValueRequest {
	return FindValueRequest{
		CallHeader: CallHeader{
			Sender:    k.routes.self,
			NetworkID: k.NetworkID,
		},
		Target: target,
	}
}

type FindValueResponse struct {
	CallHeader
	Contacts Contacts
	Value    []byte
}

func (k *Kademlia) FindValue(contact Contact, target NodeID) ([]Contact, []byte, error) {
	client, err := k.Network.Connect(contact)
	if err != nil {
		return nil, nil, err
	}

	req := k.NewFindValueRequest(target)
	res := FindValueResponse{}

	err = client.FindValue(req, &res)
	if err != nil {
		return nil, nil, err
	}

	return res.Contacts, res.Value, nil
}

func (k *Kademlia) FindValueHandler(req FindValueRequest, res *FindValueResponse) error {
	err := k.HandleCall(req.CallHeader, &res.CallHeader)
	if err != nil {
		return err
	}

	value, err := k.Storage.Get(req.Target)
	if err != nil {
		//log.Println(err)
		//panic("Read from values database failed")
		return err
	}

	if value != nil {
		res.Value = value
		return nil
	}

	res.Contacts = k.routes.FindClosest(req.Target, BucketSize)

	return nil
}
