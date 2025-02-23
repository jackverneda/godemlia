package basic

import (
	"crypto/sha256"
	"math/big"
	"strconv"
	"strings"

	"github.com/jackverneda/godemlia/pb"
)

type NodeInfo struct {
	ID   []byte `json:"id,omitempty"`
	IP   string `json:"ip,omitempty"`
	Port int    `json:"port,omitempty"`
}

func (b *NodeInfo) Equal(other NodeInfo) bool {
	return strings.EqualFold(b.IP, other.IP) &&
		b.Port == other.Port
}

func NewID(ip string, port int) ([]byte, error) {
	dumb := []byte(ip + ":" + strconv.FormatInt(int64(port), 10))
	hashValue := sha256.Sum256(dumb)
	return []byte(hashValue[:]), nil
}

func CastKBucket(nodes *[]NodeInfo) *pb.KBucket {
	result := pb.KBucket{Bucket: []*pb.NodeInfo{}}
	for _, node := range *nodes {
		result.Bucket = append(result.Bucket,
			&pb.NodeInfo{
				ID:   node.ID,
				IP:   node.IP,
				Port: int32(node.Port),
			},
		)
	}
	return &result
}

func ClosestNodeToKey(key []byte, id1 []byte, id2 []byte) int {
	buf1 := new(big.Int).SetBytes(key)
	buf2 := new(big.Int).SetBytes(id1)
	buf3 := new(big.Int).SetBytes(id2)
	result1 := new(big.Int).Xor(buf1, buf2)
	result2 := new(big.Int).Xor(buf1, buf3)

	return result1.Cmp(result2)
}
