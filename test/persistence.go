package test

import (
	"encoding/json"

	"github.com/jbenet/go-base58"
)

type Storage struct {
	KV map[string][][]byte
}

func NewStorage() *Storage {
	s := &Storage{}
	s.KV = make(map[string][][]byte)
	return s
}
func (s *Storage) Handle(action string, request string, data *[]byte) (response *[]byte, err error) {
	var payload struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if action == "CREATE" && request == "user" {
		err := json.Unmarshal(*data, &payload)
		if err != nil {
			return nil, err
		}

		err = s.Create([]byte(payload.Email), data)
		if err != nil {
			return nil, err
		}

		return nil, nil

	}
	return nil, nil
}

func (s *Storage) Create(key []byte, data *[]byte) error {
	id := base58.Encode(key)

	// _, exists := s.KV[id]
	// if !exists {
	// }

	s.KV[id] = append(s.KV[id], *data)
	//fmt.Println(s.KV[id])
	return nil
}

// func (s *Storage) Read(key []byte, start int64, end int64) (data *[]byte, err error) {
// 	id := base58.Encode(key)

// 	v, exists := s.KV[id]
// 	if !exists {
// 		return nil, errors.New("the key is not found")
// 	}
// 	////fmt.Println("Retrived v:", v)

// 	flattenByteArray := getFlattenByteArray(v)
// 	////fmt.Println("Flatten array:", flattenByteArray)
// 	return &flattenByteArray, nil
// }

// func (s *MetadataStorage) Delete(key []byte) error {
// 	id := base58.Encode(key)

// 	_, exists := s.KV[id]
// 	if !exists {
// 		return errors.New("the key is not found")
// 	}

// 	delete(s.KV, id)
// 	return nil
// }

// func (s *MetadataStorage) GetKeys() [][]byte {
// 	keys := [][]byte{}
// 	for k := range s.KV {
// 		keyStr := base58.Decode(k)
// 		keys = append(keys, keyStr)
// 	}
// 	return keys
// }

// func getFlattenByteArray(data [][]byte) []byte {
// 	flatByteArray := []byte{}

// 	for _, elem := range data {
// 		elemLen := len(elem)
// 		byteLen := make([]byte, 4)
// 		binary.LittleEndian.PutUint32(byteLen, uint32(elemLen))

// 		flatByteArray = append(flatByteArray, byteLen...)
// 		flatByteArray = append(flatByteArray, elem...)
// 	}

// 	return flatByteArray
