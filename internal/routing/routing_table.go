package routing

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	pb "github.com/jackverneda/godemlia/pb"
	godemlia "github.com/jackverneda/godemlia/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// Represents the concurrency factor in the Kademlia protocol.
	ALPHA = 3

	// The maximum number of contacts stored in a bucket, and replication parameter.
	K = 5

	// Represents the number of bits in the key space.
	B = 160
)

type RoutingTable struct {
	*godemlia.Node
	*sync.Mutex
	KBuckets [B][]godemlia.Node
}

func NewRoutingTable(b godemlia.Node) *RoutingTable {
	rt := &RoutingTable{}
	rt.ID = b.ID
	rt.IP = b.IP
	rt.Port = b.Port
	rt.Mutex = &sync.Mutex{}
	rt.KBuckets = [B][]godemlia.Node{}
	return rt
}

func (rt *RoutingTable) isAlive(b godemlia.Node) bool {
	address := fmt.Sprintf("%s:%d", rt.IP, rt.Port)
	conn, _ := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	client := pb.NewFullNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pbNode, err := client.Ping(ctx,
		&pb.Node{
			ID:   rt.ID,
			IP:   rt.IP,
			Port: int32(rt.Port),
		})

	if err != nil {
		return false
	}

	if !rt.Equal(godemlia.Node{ID: pbNode.ID, IP: pbNode.IP, Port: int(pbNode.Port)}) {
		return false
	}

	return true
}

func (rt *RoutingTable) AddNode(b godemlia.Node) error {
	rt.Lock()
	defer rt.Unlock()

	if b.Equal(*rt.Node) {
		return errors.New("ERR: Self adition into the Routing Table")
	}

	if b.Port != rt.Port {
		err := fmt.Sprintf("ERR: Invalid port %d", b.Port)
		return errors.New(err)
	}

	// get the correspondient bucket
	bIndex := getBucketIndex(rt.ID, b.ID)
	bucket := rt.KBuckets[bIndex]

	// update the node position in the case of it already
	// belongs to the bucket
	for i := 0; i < len(bucket); i++ {
		if bucket[i].Equal(b) {
			bucket = append(bucket[:i], bucket[i+1:]...)
			bucket = append(bucket, b)
			goto RETURN
		}
	}

	if len(bucket) < K {
		bucket = append(bucket, b)
	} else if !rt.isAlive(bucket[0]) {
		bucket = append(bucket[1:], b)
	}

RETURN:
	rt.KBuckets[bIndex] = bucket
	////fmt.Println(rt.KBuckets)
	return nil
}

func getBucketIndex(id1 []byte, id2 []byte) int {
	// Look at each byte from left to right
	for j := 0; j < len(id1); j++ {
		// xor the byte
		xor := id1[j] ^ id2[j]

		// check each bit on the xored result from left to right in order
		for i := 0; i < 8; i++ {
			if hasBit(xor, uint(i)) {
				byteIndex := j * 8
				bitIndex := i
				return byteIndex + bitIndex
			}
		}
	}

	// the ids must be the same
	// this should only happen during bootstrapping
	return 0
}

func hasBit(n byte, pos uint) bool {
	pos = 7 - pos
	val := n & (1 << pos)
	return (val > 0)
}

func (rt *RoutingTable) GetClosestContacts(num int, target []byte, ignoredNodes []*godemlia.Node) *ShortList {
	rt.Lock()
	defer rt.Unlock()
	// First we need to build the list of adjacent indices to our target
	// in order
	////fmt.Println(rt.NodeInfo.ID, target)
	index := getBucketIndex(rt.ID, target)
	indexList := []int{index}
	i := index - 1
	j := index + 1

	for len(indexList) < B {
		if j < B {
			indexList = append(indexList, j)
		}
		if i >= 0 {
			indexList = append(indexList, i)
		}
		i--
		j++
	}

	sl := &ShortList{Nodes: &[]godemlia.Node{}, Comparator: rt.ID}

	leftToAdd := num

	// Next we select alpha contacts and add them to the short list
	for leftToAdd > 0 && len(indexList) > 0 {
		index, indexList = indexList[0], indexList[1:]
		bucketContacts := len(rt.KBuckets[index])
		for i := 0; i < bucketContacts; i++ {
			ignored := false
			for j := 0; j < len(ignoredNodes); j++ {
				if bytes.Equal(rt.KBuckets[index][i].ID, ignoredNodes[j].ID) {
					ignored = true
				}
			}
			if !ignored {
				sl.Append([]*godemlia.Node{&rt.KBuckets[index][i]})
				leftToAdd--
				if leftToAdd == 0 {
					break
				}
			}
		}
	}

	sort.Sort(sl)

	return sl
}
