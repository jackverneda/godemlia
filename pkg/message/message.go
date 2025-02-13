package message

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"net"

	godemlia "github.com/jackverneda/godemlia/pkg"
)

type Message struct {
	IP     net.IP
	Buffer *[]byte
}

func SerializeMessage(q *[]godemlia.NodeInfo) (*[]byte, error) {
	var msgBuffer bytes.Buffer
	enc := gob.NewEncoder(&msgBuffer)

	for i := 0; i < len(*q); i++ {
		err := enc.Encode((*q)[i])
		if err != nil {
			return nil, err
		}
	}

	length := msgBuffer.Len()

	var lengthBytes [16]byte
	binary.PutUvarint(lengthBytes[:], uint64(length))

	var amountNodes [8]byte
	binary.PutUvarint(amountNodes[:], uint64(len(*q)))

	var result []byte
	result = append(result, amountNodes[:]...)
	result = append(result, lengthBytes[:]...)
	result = append(result, msgBuffer.Bytes()...)

	return &result, nil
}
