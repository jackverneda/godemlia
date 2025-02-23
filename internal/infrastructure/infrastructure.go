package infrastructure

type IInfrastructure interface {
	MainEntity() string
	Handle(action string, entity string, data *[]byte) (response *[]byte, err error)
	GetKeys() map[string][][]byte
	// Update()
	Store(entity string, id []byte, data *[]byte) (response *[]byte, err error)
	Read(entity string, id []byte) (response *[]byte, err error)
	Delete(entity string, id []byte) (err error)
}
