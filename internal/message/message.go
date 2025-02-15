package message

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"net"

	"github.com/jackverneda/godemlia/internal/basic"
)

type Message struct {
	IP     net.IP
	Buffer *[]byte
}

func SerializeMessage(q *[]basic.NodeInfo) (*[]byte, error) {
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

func DeserializeMessage(conn *net.TCPConn) (*[]basic.NodeInfo, error) {
	defer conn.Close()

	amountNodes := make([]byte, 8)
	_, err := conn.Read(amountNodes)
	if err != nil {
		return nil, err
	}

	amountReader := bytes.NewBuffer(amountNodes)
	amount, err := binary.ReadUvarint(amountReader)
	if err != nil {
		return nil, err
	}

	lengthBytes := make([]byte, 16)
	_, err = conn.Read(lengthBytes)
	if err != nil {
		return nil, err
	}

	lengthReader := bytes.NewBuffer(lengthBytes)
	length, err := binary.ReadUvarint(lengthReader)
	if err != nil {
		return nil, err
	}

	msgBytes := make([]byte, length)
	_, err = conn.Read(msgBytes)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewBuffer(msgBytes)

	resp := []basic.NodeInfo{}
	dec := gob.NewDecoder(reader)

	for i := 0; i < int(amount); i++ {
		node := basic.NodeInfo{}
		err = dec.Decode(&node)
		if err != nil {
			return nil, err
		}
		resp = append(resp, node)
	}

	return &resp, nil
}
