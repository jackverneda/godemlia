package test

import (
	"errors"
	"fmt"

	"github.com/jbenet/go-base58"
)

type Storage struct {
	KV map[string][]byte
}

func NewStorage() *Storage {
	s := &Storage{}
	s.KV = make(map[string][]byte)
	return s
}

func (s *Storage) MainEntity() string {
	return "user"
}

func (s *Storage) Handle(action string, request string, data *[]byte) (response *[]byte, err error) {
	// var payload struct {
	// 	UserName string `json:"username"`
	// 	Email    string `json:"email"`
	// 	Password string `json:"password"`
	// }

	if action == "CREATE" && request == "user" {
		// err := json.Unmarshal(*data, &payload)
		// if err != nil {
		// 	return nil, err
		// }

		// fmt.Printf("Payload: %s\n", payload.Email)
		// err = s.Create([]byte(payload.Email), data)
		// if err != nil {
		// 	return nil, err
		// }

		return nil, nil

	}
	return nil, nil
}

func (s *Storage) Store(entity string, key []byte, data *[]byte) (*[]byte, error) {
	id := base58.Encode(key)

	// _, exists := s.KV[id]
	// if !exists {
	// }

	s.KV[id] = *data
	// fmt.Println(id, s.KV[id])
	return data, nil
}

func (s *Storage) GetKeys() map[string][][]byte {
	keys := map[string][][]byte{}
	for k, _ := range s.KV {
		keys["user"] = append(keys["user"], []byte(k))
	}
	fmt.Println(keys)
	return keys
}

func (s *Storage) Read(entity string, key []byte) (data *[]byte, err error) {
	id := base58.Encode(key)

	v, exists := s.KV[id]
	if !exists {
		return nil, errors.New("the key is not found")
	}
	////fmt.Println("Retrived v:", v)

	// flattenByteArray := getFlattenByteArray(v)
	////fmt.Println("Flatten array:", flattenByteArray)
	return &v, nil
}

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
// }
