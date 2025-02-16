package infrastructure

type IInfrastructure interface {
	Handle(action string, request string, data *[]byte) (response *[]byte, err error)

	// Create()
	// Update()
	// Delete()
	// Read()
}
