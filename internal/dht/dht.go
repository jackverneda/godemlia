package dht

import (
	"bytes"

	"github.com/jackverneda/godemlia/internal/basic"
	"github.com/jackverneda/godemlia/internal/infrastructure"
	"github.com/jackverneda/godemlia/internal/routing"
)

type DHT struct {
	*routing.RoutingTable
	infrastructure.IInfrastructure
}

func (fn *DHT) Store(entity string, info []byte, data *[]byte) error {
	//fmt.Printf("INIT DHT.Store(%v) len(*data)=%d\n", key, len(*data))
	// defer //fmt.Printf("END DHT.Store(%v)\n", key)

	//fmt.Println("Before Storage.Create()")
	_, err := fn.IInfrastructure.Store(entity, info, data)
	return err
}

func (fn *DHT) FindValue(entity string, infoHash *[]byte) (value *[]byte, neighbors *[]basic.NodeInfo) {
	// value, err := fn.Handle("READ", "user", nil)
	value, err := fn.IInfrastructure.Read(entity, *infoHash)
	if err != nil {
		////fmt.Println("Find Value error: ", err)
		neighbors = fn.RoutingTable.GetClosestContacts(routing.ALPHA, *infoHash, []*basic.NodeInfo{fn.NodeInfo}).Nodes
		return nil, neighbors
	}
	return value, nil
}

func (fn *DHT) FindAll(entity, creiteria string) (value *[]byte, total int32, neighbors *[]basic.NodeInfo) {
	// value, err := fn.Handle("READ", "user", nil)
	value, total, err := fn.IInfrastructure.Search(entity, creiteria)
	if err != nil {
		////fmt.Println("Find Value error: ", err)
		neighbors = fn.RoutingTable.GetClosestContacts(routing.ALPHA, fn.ID, []*basic.NodeInfo{fn.NodeInfo}).Nodes
		return nil, 0, neighbors
	}
	return value, total, nil
}

func (fn *DHT) DeleteValue(infoHash *[]byte, start int64, end int64) error {
	_, err := fn.Handle("DELETE", "user", nil)
	if err != nil {
		//fmt.Println("ERROR line:23 DHT.Storage.Create()")
		return err
	}
	return err
}

func (fn *DHT) UpdateVaue(infoHash *[]byte, start int64, end int64) (value *[]byte, neighbors *[]basic.NodeInfo) {
	// value, err := fn.Handle(*infoHash, nil, nil)
	// if err != nil {
	// 	////fmt.Println("Find Value error: ", err)
	// 	neighbors = fn.RoutingTable.GetClosestContacts(routing.ALPHA, *infoHash, []*basic.NodeInfo{fn.NodeInfo}).Nodes
	// 	return nil, neighbors
	// }
	return value, nil
}

func (fn *DHT) FindNode(target *[]byte) (kBucket *[]basic.NodeInfo) {
	if bytes.Equal(fn.ID, *target) {
		kBucket = &[]basic.NodeInfo{*fn.NodeInfo}
	}
	kBucket = fn.RoutingTable.GetClosestContacts(routing.ALPHA, *target, []*basic.NodeInfo{fn.NodeInfo}).Nodes

	return kBucket
}
