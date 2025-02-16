package test

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

	kademlia "github.com/jackverneda/godemlia/pkg"
	base58 "github.com/jbenet/go-base58"
)

type Peer struct {
	kademlia.Node
}

func InitPeer(ip string, port, bootPort int) *Peer {
	peer := NewPeer(ip, port, bootPort, true)
	addr := fmt.Sprintf("%s:%d", ip, port)
	go peer.CreateGRPCServer(addr)
	return peer
}

func NewPeer(ip string, port, bootPort int, isBootstrapNode bool) *Peer {
	db := NewStorage()
	newPeer := kademlia.NewNode(ip, port, bootPort, db, isBootstrapNode)

	return &Peer{*newPeer}
}

func (p *Peer) Store(entity string, data *[]byte) (*[]byte, error) {

	hash := sha1.Sum(*data)
	infoHash := base58.Encode(hash[:])

	payload := map[string]interface{}{}
	err := json.Unmarshal(*data, &payload)
	if err != nil {
		return nil, err
	}

	keyHash := payload[entity+"_id"].(string)

	// fmt.Println("Before StoreValue()")
	_, err = p.StoreValue(entity, keyHash, infoHash, data)
	if err != nil {
		return nil, err
	}
	//fmt.Println("After StoreValue()")

	return data, nil
}
