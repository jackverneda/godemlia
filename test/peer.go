package test

import (
	"crypto/sha1"
	"fmt"

	kademlia "github.com/jackverneda/godemlia/pkg"
	base58 "github.com/jbenet/go-base58"
)

type Peer struct {
	kademlia.Node
}

func InitPeer(ip string, port, bootPort int) *Peer {
	peer := NewPeer(ip, port, bootPort, false)
	addr := fmt.Sprintf("%s:%d", ip, port)
	go peer.CreateGRPCServer(addr)
	return peer
}

func NewPeer(ip string, port, bootPort int, isBootstrapNode bool) *Peer {
	db := NewStorage()
	newPeer := kademlia.NewNode(ip, port, bootPort, db, isBootstrapNode)

	return &Peer{*newPeer}
}

func (p *Peer) Store(data *[]byte) (string, error) {

	hash := sha1.Sum(*data)
	key := base58.Encode(hash[:])

	//fmt.Println("Before StoreValue()")
	_, err := p.StoreValue(key, data)
	if err != nil {
		return "", nil
	}
	//fmt.Println("After StoreValue()")

	return key, nil
}
