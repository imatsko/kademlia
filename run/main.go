package main

import (
	"flag"
	"fmt"
	"github.com/imatsko/kademlia"
)

func parseFlags() (port *int, firstContact *kademlia.Contact, action bool, target string) {
	port = flag.Int("port", 6000, "a int")
	firstID := flag.String("first-id", "", "a hexideicimal node ID")
	firstIP := flag.String("first-ip", "", "the TCP address of an existing node")

	flag.BoolVar(&action, "action", false, "do some action")
	flag.StringVar(&target, "target", "", "target")

	flag.Parse()

	if *firstID == "" || *firstIP == "" {
		firstID = nil
		firstIP = nil
	} else {
		firstContact = &kademlia.Contact{}
		*firstContact = kademlia.NewContact(kademlia.NewNodeID(*firstID), *firstIP)
	}

	return
}

func main() {
	port, firstContact, action, target := parseFlags()

	if port == nil {
		panic("Must supply desired port number")
	}

	fmt.Println("Initializing Kademlia DHT ...")

	selfID := kademlia.NewRandomNodeID()

	selfAddress := fmt.Sprintf("127.0.0.1:%d", *port)
	self := kademlia.NewContact(selfID, selfAddress)
	fmt.Println("Self:", selfID, selfAddress)

	selfNetwork := kademlia.NewKademlia(self, "Certcoin-DHT")

	selfNetwork.Storage = kademlia.NewMapStorage()
	selfNetwork.Network = kademlia.NewRPCNetwork()

	nodeHandler := kademlia.NewRPCNode(selfNetwork)
	nodeHandler.Serve(selfAddress)

	if firstContact != nil {
		fmt.Println("Start bootstrap")
		contacts, err := selfNetwork.Bootstrap(*firstContact, self)
		if err != nil {
			fmt.Println("Bootstrap error:", err)
		}
		fmt.Printf("Contacts %+v \n", contacts)

		final := make(chan kademlia.Contacts)
		go selfNetwork.IterativeFindNode(firstContact.ID, kademlia.Delta, final)
		contacts = <-final
		fmt.Println("Iterative Find Node:", contacts)
	}

	if action {
		fmt.Println(target)
		fmt.Println("TEST TEST TEST")
	}

	done := make(chan bool)
	_ = <-done
}
