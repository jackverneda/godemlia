package infrastructure

type IInfrastructure interface {
	MainEntity() string
	Handle(action string, entity string, data *[]byte) (response *[]byte, err error)
	GetKeys() map[string][][]byte
	// Create()
	// Update()
	Store(entity string, key []byte, data *[]byte) (response *[]byte, err error)
	Read(entity string, key []byte) (response *[]byte, err error)
}
