kademlia/
│── cmd/                # Example CLI tool to test the library
│   ├── main.go
│── internal/           # Internal implementation details
│   ├── routing/        # Routing table & k-buckets
│   │   ├── routing_table.go
│   │   ├── kbucket.go
│   ├── network/        # gRPC server & client communication
│   │   ├── grpc_server.go
│   │   ├── grpc_client.go
│   ├── storage/        # Data storage for the DHT
│   │   ├── store.go
│── proto/              # gRPC protocol definitions
│   ├── kademlia.proto
│── pkg/                # Public API for external use
│   ├── kademlia.go     # High-level Kademlia API
│── go.mod              # Go module definition
│── README.md           # Documentation
│── tests/              # Integration & unit tests
│   ├── kademlia_test.go


## Generate ProtocolBuffer files from gRPC
export PATH="$PATH:$(go env GOPATH)/bin"
protoc --go_out=../ --go-grpc_out=../ kademlia.proto

## Test containers
docker build -t test:lastest ./test
docker run -it --network test_kademlia --network-alias test_kademlia -v "$(pwd)":/app -w /app test:lastest sh
