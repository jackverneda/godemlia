package kademlia

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sort"
	"time"

	"github.com/jackverneda/godemlia/internal/basic"
	"github.com/jackverneda/godemlia/internal/dht"
	"github.com/jackverneda/godemlia/internal/infrastructure"
	"github.com/jackverneda/godemlia/internal/message"
	"github.com/jackverneda/godemlia/internal/network"
	"github.com/jackverneda/godemlia/internal/routing"
	"github.com/jackverneda/godemlia/pb"
	"github.com/jbenet/go-base58"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Node struct {
	pb.UnimplementedNodeServer
	dht *dht.DHT
}

func NewNode(nodeIP string, nodePort, bootstrapPort int, infra infrastructure.IInfrastructure, isBootstrapNode bool) *Node {

	id, _ := basic.NewID(nodeIP, nodePort)
	node := basic.NodeInfo{ID: id, IP: nodeIP, Port: nodePort}
	dht := dht.DHT{
		RoutingTable:    routing.NewRoutingTable(node),
		IInfrastructure: infra,
	}
	Node := Node{dht: &dht}

	go Node.joinNetwork(bootstrapPort)

	if isBootstrapNode {
		go Node.bootstrap(bootstrapPort)
	}

	return &Node
}

// Create gRPC Server
func (n *Node) CreateGRPCServer(grpcServerAddress string) {
	grpcServer := grpc.NewServer()

	pb.RegisterNodeServer(grpcServer, n)
	reflection.Register(grpcServer)

	go n.Republish()

	listener, err := net.Listen("tcp", grpcServerAddress)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}
}

// ======================== RPC KADEMLIA PROTOCOL ===========================

func (n *Node) Ping(ctx context.Context, sender *pb.NodeInfo) (*pb.NodeInfo, error) {
	fmt.Printf("PING FROM %s\n\n", sender.IP)

	// add the sender to the Routing Table
	_sender := basic.NodeInfo{
		ID:   sender.ID,
		IP:   sender.IP,
		Port: int(sender.Port),
	}
	n.dht.RoutingTable.AddNode(_sender)

	receiver := &pb.NodeInfo{
		ID:   n.dht.ID,
		IP:   n.dht.IP,
		Port: int32(n.dht.Port),
	}

	return receiver, nil
}

func (n *Node) Store(ctx context.Context, data *pb.StoreData) (*pb.Response, error) {
	//fmt.Printf("INIT FullNode.Store()\n\n")
	// defer //fmt.Printf("END FullNode.Store()\n\n")

	key := data.Key
	buffer := []byte(data.Value)
	sender := basic.NodeInfo{
		ID:   data.Sender.ID,
		IP:   data.Sender.IP,
		Port: int(data.Sender.Port),
	}
	n.dht.RoutingTable.AddNode(sender)

	// //fmt.Println("Received Data:", buffer)

	err := n.dht.Store(key, &buffer)
	if err != nil {
		//fmt.Printf("ERROR line:140 DHT.Store()\n\n")
		return nil, err
	}
	return nil, nil
}

// func (n *Node) Delete(stream pb.Node_StoreServer) error {
// 	//fmt.Printf("INIT FullNode.Store()\n\n")
// 	// defer //fmt.Printf("END FullNode.Store()\n\n")

// 	key := []byte{}
// 	buffer := []byte{}
// 	var init int64 = 0

// 	for {
// 		data, err := stream.Recv()
// 		if data == nil {
// 			//fmt.Printf("END Streaming\n\n")
// 			break
// 		}
// 		if err != nil {
// 			//fmt.Printf("EXIT line:133 Store() method\n\n")
// 			return errors.New("missing chunck")
// 		}

// 		if init == 0 {
// 			//fmt.Printf("INIT Streaming\n\n")
// 			// add the sender to the Routing Table
// 			sender := basic.NodeInfo{
// 				ID:   data.Sender.ID,
// 				IP:   data.Sender.IP,
// 				Port: int(data.Sender.Port),
// 			}
// 			n.dht.RoutingTable.AddNode(sender)
// 		}

// 		key = data.Key
// 		if init == data.Value.Init {
// 			buffer = append(buffer, data.Value.Buffer...)
// 			init = data.Value.End
// 		} else {
// 			//fmt.Printf("ERROR missing chunck\n\n")
// 			return err
// 		}
// 		//fmt.Printf("OKKKK ===> FullNode(%s).Recv(%d, %d)\n", n.dht.IP, data.Value.Init, data.Value.End)
// 	}
// 	// //fmt.Println("Received Data:", buffer)

// 	err := n.dht.Store(key, &buffer)
// 	if err != nil {
// 		//fmt.Printf("ERROR line:140 DHT.Store()\n\n")
// 		return err
// 	}
// 	return nil
// }

func (n *Node) FindNode(ctx context.Context, target *pb.Target) (*pb.KBucket, error) {
	// add the sender to the Routing Table
	sender := basic.NodeInfo{
		ID:   target.Sender.ID,
		IP:   target.Sender.IP,
		Port: int(target.Sender.Port),
	}
	n.dht.RoutingTable.AddNode(sender)

	bucket := n.dht.FindNode(&target.ID)

	return basic.CastKBucket(bucket), nil
}

func (fn *Node) FindValue(ctx context.Context, target *pb.Target) (*pb.FindValueResponse, error) {
	// add the sender to the Routing Table
	sender := basic.NodeInfo{
		ID:   target.Sender.ID,
		IP:   target.Sender.IP,
		Port: int(target.Sender.Port),
	}
	fn.dht.RoutingTable.AddNode(sender)

	value, neighbors := fn.dht.FindValue(&target.ID, target.Init, target.End)
	response := pb.FindValueResponse{}

	if value == nil && neighbors != nil {
		response = pb.FindValueResponse{
			KNeartestBuckets: basic.CastKBucket(neighbors),
			Value:            []byte{},
		}

	} else if value != nil && neighbors == nil {
		//fmt.Println("Value from FindValue:", value)
		response = pb.FindValueResponse{
			KNeartestBuckets: &pb.KBucket{Bucket: []*pb.NodeInfo{}},
			Value:            *value,
		}
	} else {
		return nil, errors.New("check code because this case shouldn't be valid")
	}
	return &response, nil
}

// ======================== CORE KADEMLIA PROTOCOL ===========================

func (fn *Node) LookUp(target []byte) ([]basic.NodeInfo, error) {
	sl := fn.dht.RoutingTable.GetClosestContacts(routing.ALPHA, target, []*basic.NodeInfo{})

	contacted := make(map[string]bool)
	contacted[string(fn.dht.ID)] = true

	if len(*sl.Nodes) == 0 {
		return nil, nil
	}

	for {
		addedNodes := 0

		for i, node := range *sl.Nodes {
			if i >= routing.ALPHA {
				break
			}
			if contacted[string(node.ID)] {
				continue
			}
			contacted[string(node.ID)] = true

			// get RPC client
			client := NewNodeClient(node.IP, 8080)
			if client == nil {
				continue
			}

			// function to add the received nodes into the short list
			addRecvNodes := func(recvNodes *pb.KBucket) {
				kBucket := []*basic.NodeInfo{}
				for _, pbNode := range recvNodes.Bucket {
					if !contacted[string(pbNode.ID)] {
						kBucket = append(kBucket, &basic.NodeInfo{
							ID:   pbNode.ID,
							IP:   pbNode.IP,
							Port: int(pbNode.Port),
						})
						addedNodes++
					}
				}
				sl.Append(kBucket)
			}
			// //fmt.Println("Before timeout")
			// <-time.After(10 * time.Second)
			// //fmt.Println("After timeout")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			recvNodes, err := client.FindNode(ctx,
				&pb.Target{
					ID: node.ID,
					Sender: &pb.NodeInfo{
						ID:   fn.dht.ID,
						IP:   fn.dht.IP,
						Port: int32(fn.dht.Port),
					},
				},
			)
			if err != nil && err.Error() == "rpc error: code = DeadlineExceeded desc = context deadline exceeded" {
				//fmt.Println("Crash connection")
				sl.RemoveNode(&node)
				continue
			}
			if err != nil {
				return nil, err
			}
			addRecvNodes(recvNodes)
		}

		sl.Comparator = fn.dht.ID
		sort.Sort(sl)

		if addedNodes == 0 {
			//fmt.Println("0 added nodes")
			break
		}
	}

	kBucket := []basic.NodeInfo{}

	for i, node := range *sl.Nodes {
		if i == routing.K {
			break
		}
		//fmt.Println("append node", node.IP)
		kBucket = append(kBucket, basic.NodeInfo{
			ID:   node.ID,
			IP:   node.IP,
			Port: node.Port,
		})
	}
	return kBucket, nil
}

func (fn *Node) StoreValue(key string, data *[]byte) (string, error) {
	//fmt.Printf("INIT FullNode.StoreValue(%s) method\n\n", key)
	// defer //fmt.Printf("EXIT FullNode.StoreValue(%s) method\n\n", key)

	keyHash := base58.Decode(key)
	nearestNeighbors, err := fn.LookUp(keyHash)
	if err != nil {
		//fmt.Printf("ERROR LookUP() method\n\n")
		return "", err
	}

	if len(nearestNeighbors) < routing.K {
		err := fn.dht.Store(keyHash, data)
		if err != nil {
			fmt.Printf("ERROR DHT.Store(Me)\n\n")
		}
	}

	for index, node := range nearestNeighbors {
		if index == routing.K-1 && basic.ClosestNodeToKey(keyHash, fn.dht.ID, node.ID) == -1 {
			err := fn.dht.Store(keyHash, data)
			if err != nil {
				fmt.Printf("ERROR DHT.Store(Me)\n\n")
			}
			break
		}

		client := NewNodeClient(node.IP, node.Port)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if client == nil {
			continue
		}

		sender, err := client.Store(ctx)
		if err != nil {
			//fmt.Printf("ERROR Store(%v, %d) method", node.IP, node.Port)
			if ctx.Err() == context.DeadlineExceeded {
				// Handle timeout error
				//fmt.Println("Timeout exceeded")
				continue
			}
			//fmt.Println(err.Error())
		}
		// //fmt.Println("data bytes", dataBytes)
		err = sender.Send(
			&pb.StoreData{
				Key: keyHash,
				Sender: &pb.NodeInfo{
					ID:   fn.dht.ID,
					IP:   fn.dht.IP,
					Port: int32(fn.dht.Port),
				},
				Value: *data,
			},
		)
		if err != nil {

		}

	}

	// //fmt.Println("Stored ID: ", key, "Stored Data:", data)
	//fmt.Println("===> OKKKK")
	return key, nil
}

func (fn *Node) GetValue(target string, start int64, end int64) ([]byte, error) {
	keyHash := base58.Decode(target)

	// val, err := fn.dht.IPersistance.Read(keyHash, start, end)
	val, err := fn.dht.IInfrastructure.Handle("READ", "user", nil)
	if err == nil {
		return *val, nil
	}

	nearestNeighbors, err := fn.LookUp(keyHash)
	if err != nil {
		return nil, nil
	}
	////fmt.Println(nearestNeighbors)
	buffer := []byte{}

	for _, node := range nearestNeighbors {
		if len(target) == 0 {
			//fmt.Println("Invalid target decoding.")
			continue
		}

		clientChnn := make(chan pb.NodeClient)

		go func() {
			client := NewNodeClient(node.IP, node.Port)
			if client == nil {
				return
			}
			clientChnn <- client.NodeClient
			//fmt.Println("Channel value is: ", clientChnn)
		}()

		//fmt.Println("Init Select-Case")
		select {
		case <-time.After(5 * time.Second):
			//fmt.Println("Timeout")
			continue
		case client := <-clientChnn:
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()

			if client == nil {
				continue
			}

			//fmt.Println("Init FindValue")
			receiver, err := client.FindValue(ctx,
				&pb.Target{
					ID: keyHash,
					Sender: &pb.NodeInfo{
						ID:   fn.dht.ID,
						IP:   fn.dht.IP,
						Port: int32(fn.dht.Port),
					},
					Init: start,
					End:  end,
				},
			)
			if err != nil || receiver == nil {
				continue
			}
			//fmt.Println("End FindValue")
			if err != nil {
				if ctx.Err() == context.DeadlineExceeded {
					// Handle timeout error
					// //fmt.Println("Timeout exceeded")
					continue
				}
				//fmt.Println(err.Error())
				continue
			}
			var init int64 = 0

			for {
				data, err := receiver.Recv()
				if err != nil {
					//fmt.Println(err.Error())
					break
				}
				if data == nil {
					break
				} else if init == data.Value.Init {
					buffer = append(buffer, data.Value.Buffer...)
					init = data.Value.End
				}
			}
			////fmt.Println("Received value from STREAMING in GetValue():", buffer)
			// Received data
			if len(buffer) > 0 {
				// break
				goto RETURN
			}
		}
	}

RETURN:
	return buffer, nil
}

// ======================== JOIN NETWORK ===========================

func (fn *Node) bootstrap(port int) {
	recv := make(chan message.Message)

	broadcast := network.Broadcast{Port: port}
	go broadcast.Recv(recv)

	for {
		m := <-recv
		if m.IP.Equal(net.ParseIP(fn.dht.IP)) {
			continue
		}

		addr := net.TCPAddr{IP: m.IP, Port: port + 1}

		kBucket := fn.dht.FindNode(m.Buffer)

		*kBucket = append(*kBucket, *fn.dht.NodeInfo)

		resp, _ := message.SerializeMessage(kBucket)

		go network.Send(&addr, resp)
	}
}

func (fn *Node) joinNetwork(port int) {
	b := network.Broadcast{Port: port}
	go b.Send(&fn.dht.ID)

	recv := make(chan *net.TCPConn, 1)
	go network.Recv(recv, port+1)

	select {
	case conn := <-recv:
		kBucket, err := message.DeserializeMessage(conn)
		if err != nil {
			// fmt.Println("ERROR: Deserialize Message")
			log.Fatal(err)
		}

		if kBucket == nil {
			// fmt.Println("kBucket Received NIL")
			break
		}

		for _, node := range *kBucket {
			client := NewNodeClient(node.IP, node.Port)
			resp, _ := client.Ping(*fn.dht.NodeInfo)
			if resp.Equal(node) {
				fn.dht.RoutingTable.AddNode(node)
			}
		}
		return

	case <-time.After(5 * time.Second):
		break
	}
}

func (n *Node) PrintRoutingTable() {
	KBuckets := n.dht.RoutingTable.KBuckets

	for i := 0; i < len(KBuckets); i++ {
		for j := 0; j < len(KBuckets[i]); j++ {
			fmt.Println(KBuckets[i][j])
		}
	}
}

func (fn *Node) Republish() {
	for {
		<-time.After(time.Minute)
		// 	keys := fn.dht.IPersistance.GetKeys()
		// 	if len(keys) == 0 {
		// 		continue
		// 	}
		// 	for _, key := range keys {
		// 		data, _ := fn.dht.IPersistance.Read(key, 0, 0)
		// 		keyStr := base58.Encode(key)
		// 		fmt.Println("Key:", keyStr, "Data:", data)
		// 		if len(keyStr) == 0 || len(*data) == 0 {
		// 			break
		// 		}
		// 		go func() {
		// fn.StoreValue(keyStr, data)
		// 		}()
		// 	}
		fmt.Println("Republish worked good for node:", fn.dht.ID)
	}
}
