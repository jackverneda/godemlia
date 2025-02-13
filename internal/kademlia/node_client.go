package kademlia

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackverneda/godemlia/pb"
	godemlia "github.com/jackverneda/godemlia/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NodeClient struct {
	pb.NodeClient
	IP   string
	Port int
}

func NewNodeClient(ip string, port int) *NodeClient {
	address := fmt.Sprintf("%s:%d", ip, port)
	grpcConn := make(chan grpc.ClientConn)

	go func() {
		// stablish connection
		conn, _ := grpc.Dial(address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock())
		if conn == nil {
			return
		}
		grpcConn <- *conn
	}()

	select {
	case <-time.After(5 * time.Second):
		return nil
	case conn := <-grpcConn:
		client := pb.NewNodeClient(&conn)
		fnClient := NodeClient{
			NodeClient: client,
			IP:         ip,
			Port:       port,
		}
		return &fnClient
	}
}

func (fn *NodeClient) Ping(sender godemlia.NodeInfo) (*godemlia.NodeInfo, error) {
	nodeChnn := make(chan *pb.NodeInfo)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		node, err := fn.NodeClient.Ping(ctx,
			&pb.NodeInfo{
				ID:   sender.ID,
				IP:   sender.IP,
				Port: int32(sender.Port),
			})
		if err != nil {
			fmt.Println(err)
		}
		nodeChnn <- node
	}()

	select {
	case <-time.After(5 * time.Second):
		log := fmt.Sprintf("ERR: Node (%s:%d) doesn't respond", fn.IP, fn.Port)
		return nil, errors.New(log)
	case node := <-nodeChnn:
		return &godemlia.NodeInfo{
			ID:   node.ID,
			IP:   node.IP,
			Port: int(node.Port),
		}, nil
	}
}
