package dht

import (
	"bytes"

	"github.com/jackverneda/godemlia/internal/routing"
	"github.com/jackverneda/godemlia/internal/storage"
	godemlia "github.com/jackverneda/godemlia/pkg"
)

type DHT struct {
	*routing.RoutingTable
	storage.IPersistance
}

func (fn *DHT) Store(key []byte, data *[]byte) error {
	//fmt.Printf("INIT DHT.Store(%v) len(*data)=%d\n", key, len(*data))
	// defer //fmt.Printf("END DHT.Store(%v)\n", key)

	//fmt.Println("Before Storage.Create()")
	err := fn.IPersistance.Create(key, data)
	//fmt.Println("After Storage.Create()")
	if err != nil {
		//fmt.Println("ERROR line:23 DHT.Storage.Create()")
		return err
	}
	return nil
}

func (fn *DHT) FindValue(infoHash *[]byte, start int64, end int64) (value *[]byte, neighbors *[]godemlia.NodeInfo) {
	value, err := fn.IPersistance.Read(*infoHash, start, end)
	if err != nil {
		////fmt.Println("Find Value error: ", err)
		neighbors = fn.RoutingTable.GetClosestContacts(routing.ALPHA, *infoHash, []*godemlia.NodeInfo{fn.NodeInfo}).Nodes
		return nil, neighbors
	}
	return value, nil
}

func (fn *DHT) FindNode(target *[]byte) (kBucket *[]godemlia.NodeInfo) {
	if bytes.Equal(fn.ID, *target) {
		kBucket = &[]godemlia.NodeInfo{*fn.NodeInfo}
	}
	kBucket = fn.RoutingTable.GetClosestContacts(routing.ALPHA, *target, []*godemlia.NodeInfo{fn.NodeInfo}).Nodes

	return kBucket
}
