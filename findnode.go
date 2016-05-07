package kademlia

import (
	"container/heap"
	"sort"
)

type FindNodeRequest struct {
	CallHeader
	Target NodeID
}

func (k *Kademlia) NewFindNodeRequest(target NodeID) FindNodeRequest {
	return FindNodeRequest{
		CallHeader: CallHeader{
			Sender:    k.routes.self,
			NetworkID: k.NetworkID,
		},
		Target: target,
	}
}

type FindNodeResponse struct {
	CallHeader
	Contacts Contacts
}

func (k *Kademlia) FindNode(contact Contact, target NodeID, done chan Contacts) {
	client, err := k.Network.Connect(contact)
	if err != nil {
		done <- nil
		return
	}

	req := k.NewFindNodeRequest(target)
	res := FindNodeResponse{}

	err = client.FindNode(req, &res)
	if err != nil {
		done <- nil
		return
	}

	k.routes.Update(res.Sender)

	done <- res.Contacts
}

func (k *Kademlia) FindNodeHandler(req FindNodeRequest, res *FindNodeResponse) error {
	err := k.HandleCall(req.CallHeader, &res.CallHeader)
	if err != nil {
		return err
	}
	res.Contacts = k.routes.FindClosest(req.Target, BucketSize)
	return nil
}

func (k *Kademlia) IterativeFindNode(target NodeID, delta int, final chan Contacts) {
	done := make(chan Contacts)

	ret := make(Contacts, BucketSize)
	frontier := make(Contacts, BucketSize)
	seen := make(map[string]struct{})

	for _, node := range k.routes.FindClosest(target, delta) {
		ret = append(ret, node)
		heap.Push(&frontier, node)
		seen[node.ID.String()] = struct{}{}
	}

	pending := 0
	for i := 0; i < delta && frontier.Len() > 0; i++ {
		pending++
		contact := heap.Pop(&frontier).(Contact)
		go k.FindNode(contact, target, done)
	}

	for pending > 0 {
		nodes := <-done
		pending--
		for _, node := range nodes {
			if _, ok := seen[node.ID.String()]; !ok {
				ret = append(ret, node)
				heap.Push(&frontier, node)
				seen[node.ID.String()] = struct{}{}
			}
		}

		for pending < delta && frontier.Len() > 0 {
			pending++
			contact := heap.Pop(&frontier).(Contact)
			go k.FindNode(contact, target, done)
		}
	}

	sort.Sort(ret)
	if ret.Len() > BucketSize {
		ret = ret[:BucketSize]
	}

	final <- ret
}
